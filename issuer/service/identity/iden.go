package identity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"issuer/service/claim"
	"issuer/service/command"
	"issuer/service/communication"
	issuer_contract "issuer/service/contract"
	"issuer/service/identity/state"
	"issuer/service/schema"
	"math/big"
)

type Identity struct {
	sk          babyjub.PrivateKey
	Identifier  *core.ID
	authClaimId *uuid.UUID
	baseUrl     string

	state       *state.IdentityState
	CmdHandler  *command.Handler
	CommHandler *communication.Handler
}

func New(s *state.IdentityState, cmdHandler *command.Handler, commHandler *communication.Handler, sk babyjub.PrivateKey, hostUrl string) (*Identity, error) {
	iden := &Identity{
		state:       s,
		CmdHandler:  cmdHandler,
		CommHandler: commHandler,

		sk:      sk,
		baseUrl: hostUrl,
	}

	id, authClaimId, err := iden.state.GetIdentityFromDB()
	if err != nil {
		return nil, fmt.Errorf("error on identitiy initialization, %v", err)
	}

	if len(id) > 0 { // case: identity found -> load identity
		iden.Identifier = id
		iden.authClaimId = authClaimId
		return iden, nil
	}

	return iden, iden.init()
}

func (i *Identity) init() error {
	logger.Debug("Identity.init() invoked")

	logger.Debug("setup genesis state")
	identifier, authClaim, err := i.state.SetupGenesisState(i.sk.Public())
	if err != nil {
		return err
	}

	i.Identifier = identifier
	logger.Debugf("identity identifier: %v", i.Identifier)

	logger.Debug("generating auth claim proof")

	proof, err := i.generateProof(authClaim)
	if err != nil {
		return err
	}

	authClaimModel, err := claim.CoreClaimToClaimModel(authClaim, schema.AuthBJJCredentialURL, schema.AuthBJJCredential)
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
	authClaimModel.MTPProof = proof
	authClaimModel.Identifier = i.Identifier.String()
	authClaimModel.ID = uuid.New()

	logger.Debugf("adding auth claim to db, claim-id: %x", authClaimModel.ID.String())
	err = i.state.AddClaimToDB(authClaimModel)
	if err != nil {
		return err
	}
	i.authClaimId = &authClaimModel.ID

	return i.state.SaveIdentity(identifier, authClaimModel.ID)
}

func (i *Identity) generateProof(claim *core.Claim) ([]byte, error) {
	hIndex, err := claim.HIndex()
	if err != nil {
		return nil, err
	}

	proof, _, err := i.state.Claims.Tree.GenerateProof(context.Background(), hIndex, nil)
	if err != nil {
		return nil, err
	}

	mtProof := verifiable.Iden3SparseMerkleProof{}
	mtProof.Type = verifiable.Iden3SparseMerkleProofType
	mtProof.MTP = proof

	stateHash, err := i.state.GetStateHash()
	stateHashHex := stateHash.Hex()
	claimsRootHex := i.state.Claims.Tree.Root().Hex()
	mtProof.IssuerData = verifiable.IssuerData{
		ID: i.Identifier,
		State: verifiable.State{
			Value:          &stateHashHex,
			ClaimsTreeRoot: &claimsRootHex,
		},
		MTP: proof,
	}

	proofB, err := json.Marshal(mtProof)
	if err != nil {
		return nil, err
	}

	return proofB, nil
}

func (i *Identity) CreateClaim(cReq *issuer_contract.CreateClaimRequest) (*issuer_contract.CreateClaimResponse, error) {
	logger.Debug("CreateClaim() invoked")

	logger.Tracef("process schema - url: %s", cReq.Schema.URL)
	slots, encodedSchema, err := schema.Process(cReq.Schema.URL, cReq.Schema.Type, cReq.Data)
	if err != nil {
		return nil, err
	}

	claimReq := &claim.CoreClaimData{
		EncodedSchema:   encodedSchema,
		Slots:           *slots,
		SubjectID:       cReq.Identifier,
		Expiration:      cReq.Expiration,
		Version:         cReq.Version,
		Nonce:           cReq.RevNonce,
		SubjectPosition: cReq.SubjectPosition,
	}

	logger.Trace("generating core-claim from the request")
	coreClaim, err := claim.GenerateCoreClaim(claimReq)
	if err != nil {
		return nil, err
	}

	claimModel, err := claim.CoreClaimToClaimModel(coreClaim, cReq.Schema.URL, cReq.Schema.Type)
	if err != nil {
		return nil, err
	}

	// set credential status
	issuerIDString := i.Identifier.String()
	cs, err := claim.CreateCredentialStatus(i.baseUrl, verifiable.SparseMerkleTreeProof, claimModel.RevNonce)
	if err != nil {
		return nil, err
	}

	claimModel.CredentialStatus = cs
	logger.Trace("finished creating the claim object from the user request")

	logger.Trace("signing claim entry")
	newClaimSig, err := claim.SignClaimEntry(coreClaim, i.sign)
	if err != nil {
		return nil, err
	}

	authClaim, err := i.state.Claims.GetClaim([]byte(i.authClaimId.String()))
	if err != nil {
		return nil, err
	}

	sigProof, err := claim.ConstructSigProof(authClaim, newClaimSig)
	if err != nil {
		return nil, err
	}

	sigProof.IssuerData.RevocationStatus = fmt.Sprintf("%s/api/v1/claims/revocation/status/%d", i.baseUrl, claimModel.RevNonce)

	//mtProof, err := i.generateProof(authClaim.CoreClaim)
	//mtProof, err := i.generateProof(coreClaim)
	//if err != nil {
	//	return nil, err
	//}

	// Save
	claimModel.Identifier = issuerIDString
	claimModel.Issuer = issuerIDString
	claimModel.ID = uuid.New()
	jsonSignatureProof, err := json.Marshal(sigProof)
	if err != nil {
		return nil, err
	}
	claimModel.SignatureProof = jsonSignatureProof
	//claimModel.MTPProof = mtProof
	claimModel.Data = cReq.Data

	//logger.Trace("adding claim to the claims tree")
	//err = i.state.AddClaimToTree(coreClaim)
	//if err != nil {
	//	return nil, err
	//}

	logger.Trace("adding claim to the claims DB")
	err = i.state.AddClaimToDB(claimModel)
	if err != nil {
		return nil, err
	}

	return &issuer_contract.CreateClaimResponse{ID: claimModel.ID.String()}, nil
}

func (i *Identity) GetClaim(id string) (*issuer_contract.GetClaimResponse, error) {
	logger.Debug("GetClaim() invoked")

	claimID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	claimModel, err := i.state.Claims.GetClaim([]byte(claimID.String()))
	if err != nil {
		return nil, err
	}

	c, err := claim.ClaimModelToIden3Credential(claimModel)
	if err != nil {
		return nil, err
	}

	res := issuer_contract.GetClaimResponse(c)

	return &res, nil
}

func (i *Identity) GetIdentity() (*issuer_contract.GetIdentityResponse, error) {
	logger.Debug("GetIdentity() invoked")

	stateHash, err := i.state.GetStateHash()
	if err != nil {
		return nil, err
	}

	res := &issuer_contract.GetIdentityResponse{
		Identifier: i.Identifier.String(),
		State: &issuer_contract.IdentityState{
			Identifier:         i.Identifier.String(),
			State:              stateHash.Hex(),
			ClaimsTreeRoot:     i.state.Claims.Tree.Root().Hex(),
			RevocationTreeRoot: i.state.Revocations.Tree.Root().Hex(),
			RootOfRoots:        i.state.Roots.Tree.Root().Hex(),
		},
	}

	return res, nil
}

func (i *Identity) GetRevocationStatus(nonce uint64) (*issuer_contract.GetRevocationStatusResponse, error) {
	logger.Debug("GetRevocationStatus() invoked")

	rID := new(big.Int).SetUint64(nonce)

	res := &issuer_contract.GetRevocationStatusResponse{}
	mtp, err := i.state.Revocations.GenerateRevocationProof(rID)
	if err != nil {
		return nil, err
	}
	res.MTP = mtp
	res.Issuer.RevocationTreeRoot = i.state.Revocations.Tree.Root().Hex()
	res.Issuer.RootOfRoots = i.state.Roots.Tree.Root().Hex()
	res.Issuer.ClaimsTreeRoot = i.state.Claims.Tree.Root().Hex()

	stateHash, err := i.state.GetStateHash()
	if err != nil {
		return nil, err
	}
	res.Issuer.State = stateHash.Hex()

	return res, nil

}

// data should be a little-endian bytes representation of *big.Int.
//func (i *Identity) signBytes(data []byte) ([]byte, error) {
//	if len(data) > 32 {
//		return nil, errors.New("data to signBytes is too large")
//	}
//
//	z := new(big.Int).SetBytes(utils.SwapEndianness(data))
//
//	return i.sign(z)
//}

func (i *Identity) sign(z *big.Int) ([]byte, error) {
	if !utils.CheckBigIntInField(z) {
		return nil, errors.New("data to signBytes is too large")
	}

	sig := i.sk.SignPoseidon(z).Compress()
	return sig[:], nil
}
