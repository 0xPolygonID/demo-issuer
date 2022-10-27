package identity

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	issuer_contract "issuer/service/contract"
	"issuer/service/identity/state"
	"issuer/service/models"
	"issuer/service/schema"
	util "issuer/service/utils"
	"math/big"
	"strings"
	"time"
)

type Identity struct {
	sk          babyjub.PrivateKey
	Identifier  *core.ID
	authClaimId []byte
	state       *state.IdentityState
	baseUrl     string
}

var (
	BabyJubSignatureType = "BJJSignature2021"
)

func New(s *state.IdentityState, sk babyjub.PrivateKey, hostUrl string) (*Identity, error) {
	iden := &Identity{
		state:   s,
		sk:      sk,
		baseUrl: hostUrl,
	}

	err := iden.init()
	if err != nil {
		return nil, fmt.Errorf("error on identitiy initialization, %v", err)
	}

	return iden, nil
}

func (i *Identity) init() error {
	schemaHash, err := core.NewSchemaHashFromHex(models.AuthBJJCredentialHash)
	if err != nil {
		return err
	}

	authClaim, err := newAuthClaim(i.sk.Public(), schemaHash)
	if err != nil {
		return err
	}

	authClaim.SetRevocationNonce(0)

	index, value, err := authClaim.HiHv()
	if err != nil {
		return err
	}

	err = i.state.Claims.ClaimTree.Add(context.Background(), index, value)
	if err != nil {
		return err
	}

	currState, err := i.state.GetStateHash()
	if err != nil {
		return err
	}

	identifier, err := core.IdGenesisFromIdenState(core.TypeDefault, currState.BigInt())
	if err != nil {
		return err
	}

	i.Identifier = identifier

	authClaimModel, err := claimToClaimModel(authClaim, models.AuthBJJCredentialURL, models.AuthBJJCredential)
	if err != nil {
		return err
	}

	claimData := make(map[string]interface{})
	claimData["x"] = i.sk.Public().X.String()
	claimData["y"] = i.sk.Public().Y.String()
	marshalledClaimData, err := json.Marshal(claimData)
	if err != nil {
		return err
	}

	authClaimModel.Data = marshalledClaimData
	authClaimModel.Issuer = i.Identifier.String()

	proof, _, err := i.state.Claims.ClaimTree.GenerateProof(context.Background(), index, nil)
	if err != nil {
		return err
	}

	mtProof := verifiable.Iden3SparseMerkleProof{}
	mtProof.Type = verifiable.Iden3SparseMerkleProofType
	mtProof.MTP = proof

	stateHex := currState.Hex()
	claimsRootHex := i.state.Claims.ClaimTree.Root().Hex()
	mtProof.IssuerData = verifiable.IssuerData{
		ID: i.Identifier,
		State: verifiable.State{
			Value:          &stateHex,
			ClaimsTreeRoot: &claimsRootHex,
		},
	}

	proofB, err := json.Marshal(mtProof)
	if err != nil {
		return err
	}

	strIdentifier := i.Identifier.String()
	authClaimModel.ID = index.Bytes()
	authClaimModel.MTPProof = proofB
	authClaimModel.Identifier = &strIdentifier

	err = i.state.Claims.SaveClaimDB(authClaimModel)
	if err != nil {
		return err
	}
	i.authClaimId = authClaimModel.ID

	return err
}

func (i *Identity) getAuthClaim() (*models.Claim, error) {
	claim, err := i.state.Claims.GetClaim(i.authClaimId)
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func claimToClaimModel(claim *core.Claim, schemaURL, schemaType string) (*models.Claim, error) {
	otherIdentifier := ""
	id, err := claim.GetID()
	if err != nil {

	}
	otherIdentifier = id.String()

	switch err {
	case nil:
		otherIdentifier = id.String()
	case core.ErrNoID:
		return nil, fmt.Errorf("can't get ID, %v", err)
	default:
		return nil, fmt.Errorf("encounter error on claim-id extraction, %v", err)
	}

	var expiration int64
	if expirationDate, ok := claim.GetExpirationDate(); ok {
		expiration = expirationDate.Unix()
	}

	hindex, err := claim.HIndex()
	if err != nil {
		return nil, err
	}

	sb := claim.GetSchemaHash()
	schemaHash := hex.EncodeToString(sb[:])
	res := models.Claim{
		SchemaHash:      schemaHash,
		SchemaURL:       schemaURL,
		SchemaType:      schemaType,
		OtherIdentifier: otherIdentifier,
		Expiration:      expiration,
		Updatable:       claim.GetFlagUpdatable(),
		Version:         claim.GetVersion(),
		RevNonce:        claim.GetRevocationNonce(),
		CoreClaim:       claim,
		HIndex:          hindex.String(),
	}

	return &res, nil
}

func newAuthClaim(key *babyjub.PublicKey, schemaHash core.SchemaHash) (*core.Claim, error) {
	revNonce, err := util.Rand()
	if err != nil {
		return nil, err
	}

	return core.NewClaim(schemaHash,
		core.WithIndexDataInts(key.X, key.Y),
		core.WithRevocationNonce(revNonce))
}

func (i *Identity) CreateClaim(cReq *issuer_contract.CreateClaimRequest) (*issuer_contract.CreateClaimResponse, error) {

	schemaBytes, _, err := schema.Load(cReq.Schema.URL)
	if err != nil {
		return nil, err
	}

	slots, err := schema.Process(cReq.Schema.URL, cReq.Schema.Type, cReq.Data)
	if err != nil {
		return nil, err
	}

	claimReq := ClaimRequest{
		EncodedSchema:   schema.CreateSchemaHash(schemaBytes, cReq.Schema.Type),
		Slots:           slots,
		SubjectID:       cReq.Identifier,
		Expiration:      cReq.Expiration,
		Version:         cReq.Version,
		Nonce:           cReq.RevNonce,
		SubjectPosition: cReq.SubjectPosition,
	}

	coreClaim, err := GenerateCoreClaim(claimReq)
	if err != nil {
		return nil, err
	}

	claimModel, err := claimToClaimModel(coreClaim, cReq.Schema.URL, cReq.Schema.Type)
	if err != nil {
		return nil, err
	}

	// set credential status
	issuerIDString := i.Identifier.String()
	cs, err := createCredentialStatus(i.baseUrl, verifiable.SparseMerkleTreeProof, claimModel.RevNonce)
	if err != nil {
		return nil, err
	}
	claimModel.CredentialStatus = cs

	authClaim, err := i.getAuthClaim()
	if err != nil {
		return nil, err
	}

	proof, err := i.SignClaimEntry(authClaim, coreClaim)
	if err != nil {
		return nil, err
	}

	proof.IssuerData.RevocationStatus = fmt.Sprintf("%s/api/v1/claims/revocation/status/%d", i.Identifier, claimModel.RevNonce)

	// Save
	claimModel.Identifier = &issuerIDString
	claimModel.Issuer = issuerIDString

	jsonSignatureProof, err := json.Marshal(proof)
	if err != nil {
		return nil, err
	}
	claimModel.SignatureProof = jsonSignatureProof
	claimModel.Data = cReq.Data

	err = i.state.Claims.SaveClaimDB(claimModel)
	if err != nil {
		return nil, err
	}

	return &issuer_contract.CreateClaimResponse{ID: string(claimModel.ID)}, nil
}

func (i *Identity) SignClaimEntry(authClaim *models.Claim, claimEntry *core.Claim) (*verifiable.BJJSignatureProof2021, error) {
	hashIndex, hashValue, err := claimEntry.HiHv()
	if err != nil {
		return nil, err
	}

	commonHash, err := poseidon.Hash([]*big.Int{hashIndex, hashValue})
	if err != nil {
		return nil, err
	}

	issuerMTP := verifiable.Iden3SparseMerkleProof{}
	err = json.Unmarshal(authClaim.MTPProof, issuerMTP)
	if err != nil {
		return nil, err
	}

	sig, err := i.sign(commonHash)
	if err != nil {
		return nil, err
	}

	// followed https://w3c-ccg.github.io/ld-proofs/
	proof := verifiable.BJJSignatureProof2021{}
	proof.Type = BabyJubSignatureType
	proof.Signature = hex.EncodeToString(sig)
	issuerMTP.IssuerData.AuthClaim = authClaim.CoreClaim
	proof.IssuerData = issuerMTP.IssuerData

	return &proof, nil
}

// data should be a little-endian bytes representation of *big.Int.
func (i *Identity) signBytes(data []byte) ([]byte, error) {
	if len(data) > 32 {
		return nil, errors.New("data to signBytes is too large")
	}

	z := new(big.Int).SetBytes(utils.SwapEndianness(data))

	return i.sign(z)
}

func (i *Identity) sign(z *big.Int) ([]byte, error) {
	if !utils.CheckBigIntInField(z) {
		return nil, errors.New("data to signBytes is too large")
	}

	sig := i.sk.SignPoseidon(z).Compress()
	return sig[:], nil
}

func createCredentialStatus(urlBase string, sType verifiable.CredentialStatusType, revNonce uint64) ([]byte, error) {
	cStatus := verifiable.CredentialStatus{
		ID:   fmt.Sprintf("%s/api/v1/claims/revocation/status/%d", urlBase, revNonce),
		Type: sType,
	}

	b, err := json.Marshal(cStatus)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ClaimModelToIden3Credential(c *models.Claim) (*verifiable.Iden3Credential, error) {
	claimIdPos, err := getClaimIdPosition(c.CoreClaim)
	if err != nil {
		return nil, err
	}

	credSubjects := make(map[string]interface{})
	err = json.Unmarshal(c.Data, &credSubjects)
	if err != nil {
		return nil, err
	}
	credSubjects["type"] = c.SchemaType
	if len(c.OtherIdentifier) > 0 {
		credSubjects["id"] = c.OtherIdentifier
	}

	// * create proof object
	proofs := make([]interface{}, 0)

	signatureProof := &verifiable.BJJSignatureProof2021{}
	if c.SignatureProof != nil && string(c.CredentialStatus) != "{}" {
		err = json.Unmarshal(c.SignatureProof, signatureProof)
		if err != nil {
			return nil, err
		}

		proofs = append(proofs, signatureProof)
	}

	mtpProof := &verifiable.Iden3SparseMerkleProof{}
	if !strings.EqualFold(string(c.MTPProof), "{}") {
		err = json.Unmarshal(c.MTPProof, mtpProof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, mtpProof)
	}

	// create credential status object
	credStatus := &verifiable.CredentialStatus{}
	if c.CredentialStatus != nil && string(c.CredentialStatus) != "{}" {
		err = json.Unmarshal(c.CredentialStatus, credStatus)
		if err != nil {
			return nil, err
		}
	}

	res := &verifiable.Iden3Credential{
		ID:         string(c.ID),
		Type:       []string{schema.Iden3CredentialSchema},
		Expiration: time.Unix(c.Expiration, 0),
		RevNonce:   c.RevNonce,
		Updatable:  c.Updatable,
		Version:    c.Version,
		Context:    []string{schema.Iden3CredentialSchemaURL, c.SchemaURL},
		CredentialSchema: struct {
			ID   string `json:"@id"`
			Type string `json:"type"`
		}{
			ID:   c.SchemaURL,
			Type: c.SchemaType,
		},
		SubjectPosition:   claimIdPos,
		CredentialSubject: credSubjects,
		Proof:             proofs,
		CredentialStatus:  credStatus,
	}

	return res, nil
}

func (i *Identity) GetClaim(id string) (*issuer_contract.GetClaimResponse, error) {
	claimModel, err := i.state.Claims.GetClaim([]byte(id))
	if err != nil {
		return nil, err
	}

	c, err := ClaimModelToIden3Credential(claimModel)

	res := issuer_contract.GetClaimResponse(c)

	return &res, nil
}

func getClaimIdPosition(c *core.Claim) (string, error) {
	claimIdPos, err := c.GetIDPosition()
	if err != nil {
		return "", err
	}

	return subjectPositionIDToString(claimIdPos)
}

func subjectPositionIDToString(p core.IDPosition) (string, error) {
	switch p {
	case core.IDPositionIndex:
		return SubjectPositionIndex, nil
	case core.IDPositionValue:
		return SubjectPositionValue, nil
	default:
		return "", fmt.Errorf("id position is not specified")
	}
}

func (i *Identity) GetIdentity() (*issuer_contract.GetIdentityResponse, error) {
	stateHash, err := i.state.GetStateHash()
	if err != nil {
		return nil, err
	}

	res := &issuer_contract.GetIdentityResponse{
		Identifier: i.Identifier.String(),
		State: &issuer_contract.IdentityState{
			Identifier:         i.Identifier.String(),
			State:              stateHash.Hex(),
			ClaimsTreeRoot:     i.state.Claims.ClaimTree.Root().Hex(),
			RevocationTreeRoot: i.state.Revocations.RevTree.Root().Hex(),
			RootOfRoots:        i.state.Roots.RootsTree.Root().Hex(),
		},
	}

	return res, nil
}

func (i *Identity) GetRevocationStatus(nonce uint64) (*issuer_contract.GetRevocationStatusResponse, error) {
	rID := new(big.Int).SetUint64(nonce)

	res := &issuer_contract.GetRevocationStatusResponse{}
	mtp, err := i.state.Revocations.GenerateRevocationProof(rID)
	if err != nil {
		return nil, err
	}
	res.MTP = mtp
	res.Issuer.RevocationTreeRoot = i.state.Revocations.RevTree.Root().Hex()
	res.Issuer.RootOfRoots = i.state.Roots.RootsTree.Root().Hex()
	res.Issuer.ClaimsTreeRoot = i.state.Claims.ClaimTree.Root().Hex()

	stateHash, err := i.state.GetStateHash()
	if err != nil {
		return nil, err
	}
	res.Issuer.State = stateHash.Hex()

	return res, nil

}
