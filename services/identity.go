package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	store "github.com/demonsh/smt-bolt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	"lightissuer/models"
	commonutil "lightissuer/utils"
	"math/big"
)

var (
	BabyJubSignatureType = "BJJSignature2021"
)

// IdentityService service
type IdentityService struct {
	privateKey    *babyjub.PrivateKey
	publicKey     *babyjub.PublicKey
	dbService     *DBService
	treeStorage   *store.BoltStore
	schemaService *SchemaService
	MTS           *IdentityMerkleTrees
	Identity      *models.Identity
}

func NewIdentityService(ctx context.Context, privateKey *babyjub.PrivateKey, dbService *DBService, schemaService *SchemaService) (*IdentityService, error) {

	treeStorage, err := store.NewBoltStorage(dbService.Storage)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()

	iService := IdentityService{privateKey, publicKey, dbService, treeStorage, schemaService, &IdentityMerkleTrees{}, &models.Identity{}}

	_, err = iService.LoadIdentity(ctx)
	if err == ErrKeyNotFound {

		_, err = iService.NewIdentity(ctx)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	fmt.Println("identifier ->", iService.Identity.Identifier)

	return &iService, nil
}

func (i *IdentityService) NewIdentity(ctx context.Context) (*models.Identity, error) {

	dbService := i.dbService
	mts, err := CreateIdentityMerkleTrees(ctx, i.treeStorage)
	if err != nil {
		return nil, err
	}

	schemaHash, err := core.NewSchemaHashFromHex(models.AuthBJJCredentialHash)
	if err != nil {
		return nil, err
	}

	authClaim, err := newAuthClaim(i.publicKey, schemaHash)
	if err != nil {
		return nil, err
	}

	authClaim.SetRevocationNonce(0)

	index, value, err := authClaim.HiHv()
	if err != nil {
		return nil, err
	}

	err = mts.ClaimTree.Add(ctx, index, value)
	if err != nil {
		return nil, err
	}

	currentState, err := merkletree.HashElems(mts.ClaimTree.Root().BigInt(),
		merkletree.HashZero.BigInt(), merkletree.HashZero.BigInt())
	if err != nil {
		return nil, err
	}

	identifier, err := core.IdGenesisFromIdenState(core.TypeDefault, currentState.BigInt())
	if err != nil {
		return nil, err
	}

	identity := models.NewIdentityFromIdentifier(identifier,
		currentState.Hex())

	claimsTreeHex := mts.ClaimTree.Root().Hex()

	hashZeroString := "0000000000000000000000000000000000000000000000000000000000000000"

	identity.State.ClaimsTreeRoot = &claimsTreeHex
	identity.State.RevocationTreeRoot = &hashZeroString
	identity.State.RootOfRoots = &hashZeroString

	claimData := make(map[string]interface{})
	claimData["x"] = i.publicKey.X.String()
	claimData["y"] = i.publicKey.Y.String()

	marshalledClaimData, err := json.Marshal(claimData)
	if err != nil {
		return nil, err
	}

	authClaimModel, err := models.FromClaimer(authClaim, models.AuthBJJCredentialURL, models.AuthBJJCredential)
	if err != nil {
		return nil, err
	}

	authClaimModel.CoreClaim = authClaim

	proof, _, err := mts.ClaimTree.GenerateProof(ctx, index, nil)
	if err != nil {
		return nil, err
	}

	authClaimModel.Data = string(marshalledClaimData)

	stateHex := currentState.Hex()
	cltHex := mts.ClaimTree.Root().Hex()

	mtpProof := verifiable.Iden3SparseMerkleProof{
		Type: verifiable.Iden3SparseMerkleProofType,
		IssuerData: verifiable.IssuerData{
			ID: identifier,
			State: verifiable.State{
				ClaimsTreeRoot: &cltHex,
				Value:          &stateHex,
			},
			MTP: proof,
		},
	}

	jsonProof, err := json.Marshal(mtpProof)

	if err != nil {
		return nil, err
	}

	authClaimModel.MTPProof = jsonProof
	authClaimModel.Issuer = identifier.String()

	authClaimModel.IdentityState = identity.State.State
	authClaimModel.Identifier = &identity.Identifier

	_, err = dbService.saveClaim(authClaimModel)
	if err != nil {
		return nil, err
	}

	i.MTS = mts
	i.Identity = identity

	err = i.dbService.saveIdentity(i.Identity)

	if err != nil {
		return nil, err
	}

	fmt.Println("public key is ->", i.publicKey.String())
	fmt.Println("identifier ->", identity.Identifier)

	return identity, nil
}

func (i *IdentityService) LoadIdentity(_ context.Context) (*models.Identity, error) {
	identity, err := i.dbService.getIdentity()
	if err != nil {
		return nil, err
	}

	mts, err := CreateIdentityMerkleTrees(context.Background(), i.treeStorage)
	if err != nil {
		return nil, err
	}

	i.Identity = identity
	i.MTS = mts

	return identity, nil
}

// SignClaimEntry signs claim entry with provided authClaim
func (i *IdentityService) SignClaimEntry(authClaim *models.Claim, claimEntry *core.Claim) (*verifiable.BJJSignatureProof2021, error) {

	hashIndex, hashValue, err := claimEntry.HiHv()

	if err != nil {
		return nil, err
	}

	commonHash, err := poseidon.Hash([]*big.Int{hashIndex, hashValue})

	if err != nil {
		return nil, err
	}

	var issuerMTP verifiable.Iden3SparseMerkleProof
	err = json.Unmarshal(authClaim.MTPProof, &issuerMTP)
	if err != nil {
		return nil, err
	}

	signtureBytes, err := i.Sign(context.Background(), merkletree.SwapEndianness(commonHash.Bytes()))

	if err != nil {
		return nil, err
	}

	// followed https://w3c-ccg.github.io/ld-proofs/
	var proof verifiable.BJJSignatureProof2021
	proof.Type = BabyJubSignatureType
	proof.Signature = hex.EncodeToString(signtureBytes)
	issuerMTP.IssuerData.AuthClaim = authClaim.CoreClaim
	proof.IssuerData = issuerMTP.IssuerData

	return &proof, nil
}

// newAuthClaim generate BabyJubKeyTypeAuthorizeKSign claimL
func newAuthClaim(key *babyjub.PublicKey, schemaHash core.SchemaHash) (*core.Claim, error) {
	revNonce, err := commonutil.RandInt64()
	if err != nil {
		return nil, err
	}

	return core.NewClaim(schemaHash,
		core.WithIndexDataInts(key.X, key.Y),
		core.WithRevocationNonce(revNonce))
}

// Sign signs *big.Int using poseidon algorithm.
// data should be a little-endian bytes representation of *big.Int.
func (i *IdentityService) Sign(_ context.Context, data []byte) ([]byte, error) {

	if len(data) > 32 {
		return nil, errors.New("data to sign is too large")
	}

	z := new(big.Int).SetBytes(utils.SwapEndianness(data))
	if !utils.CheckBigIntInField(z) {
		return nil, errors.New("data to sign is too large")
	}

	privKey := i.privateKey

	sig := privKey.SignPoseidon(z).Compress()
	return sig[:], nil
}
