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
	"issuer/service/cfgs"
	"issuer/service/claim"
	"issuer/service/command"
	"issuer/service/communication"
	"issuer/service/identity/state"
	issuer_contract "issuer/service/models"
	"issuer/service/schema"
	"math/big"
)

type Identity struct {
	sk           babyjub.PrivateKey
	Identifier   *core.ID
	authClaimId  *uuid.UUID
	authClaim    *core.Claim
	publicUrl    string
	circuitsPath string

	state         *state.IdentityState
	CmdHandler    *command.Handler
	CommHandler   *communication.Handler
	schemaBuilder *schema.Builder
	stateStore    StateStore
}

func New(
	s *state.IdentityState,
	schemaBuilder *schema.Builder,
	sk babyjub.PrivateKey,
	cfg *cfgs.IssuerConfig,
	stateStore StateStore,
) (*Identity, error) {
	iden := &Identity{
		state:         s,
		schemaBuilder: schemaBuilder,

		sk:           sk,
		publicUrl:    cfg.PublicUrl,
		circuitsPath: cfg.KeyDir,
		stateStore:   stateStore,
	}

	id, authClaimId, err := iden.state.GetIdentityFromDB()
	if err != nil {
		return nil, fmt.Errorf("error on identitiy initialization, %v", err)
	}

	if id != nil && len(id) > 0 { // case: identity found -> load identity
		iden.Identifier = id
		iden.authClaimId = authClaimId
		iden.state.CommittedState = state.CommittedState{
			IsLatestStateGenesis: iden.state.IsGenesis(),
			RootsTreeRoot:        iden.state.Roots.Tree.Root(),
			ClaimsTreeRoot:       iden.state.Claims.Tree.Root(),
			RevocationTreeRoot:   iden.state.Revocations.Tree.Root(),
		}
		ac, err := iden.state.Claims.GetClaim([]byte(authClaimId.String()))
		if err != nil {
			return nil, err
		}
		iden.authClaim = ac.CoreClaim
	} else { // case: identity not found -> init new identity
		err = iden.init()
		if err != nil {
			return nil, fmt.Errorf("error on identitiy initialization, %v", err)
		}
	}

	iden.CommHandler = communication.NewCommunicationHandler(iden.Identifier.String(), *cfg)
	iden.CmdHandler = command.NewHandler(iden.state, cfg.KeyDir)

	return iden, nil
}

func (i *Identity) init() error {
	logger.Debug("Identity.init() invoked")

	logger.Debug("setup genesis state")
	identifier, authClaim, err := i.state.SetupGenesisState(i.sk.Public())
	if err != nil {
		return err
	}

	i.authClaim = authClaim
	i.state.CommittedState = state.CommittedState{
		IsLatestStateGenesis: true,
		ClaimsTreeRoot:       i.state.Claims.Tree.Root(),
		RevocationTreeRoot:   i.state.Revocations.Tree.Root(),
		RootsTreeRoot:        i.state.Roots.Tree.Root(),
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
	slots, encodedSchema, err := i.schemaBuilder.Process(cReq.Schema.URL, cReq.Schema.Type, cReq.Data)
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

	hi, hv, err := coreClaim.HiHv()
	if err != nil {
		return nil, err
	}

	err = i.state.Claims.Tree.Add(context.TODO(), hi, hv)
	if err != nil {
		return nil, err
	}

	claimModel, err := claim.CoreClaimToClaimModel(coreClaim, cReq.Schema.URL, cReq.Schema.Type)
	if err != nil {
		return nil, err
	}

	// set credential status
	issuerIDString := i.Identifier.String()
	cs, err := claim.CreateCredentialStatus(i.publicUrl, verifiable.SparseMerkleTreeProof, claimModel.RevNonce)
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

	sigProof.IssuerData.RevocationStatus = fmt.Sprintf("%s/api/v1/claims/revocations/%d", i.publicUrl, authClaim.RevNonce)

	// Save
	claimModel.Identifier = issuerIDString
	claimModel.Issuer = issuerIDString
	claimModel.ID = uuid.New()
	jsonSignatureProof, err := json.Marshal(sigProof)
	if err != nil {
		return nil, err
	}
	claimModel.SignatureProof = jsonSignatureProof
	claimModel.Data = cReq.Data

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

	if !i.state.CommittedState.IsLatestStateGenesis {
		claimIdx, err := claimModel.CoreClaim.HIndex()
		if err != nil {
			return nil, err
		}
		mtp, err := i.getMTPProof(claimIdx)
		if err != nil {
			return nil, err
		}
		claimModel.MTPProof, err = json.Marshal(mtp)
		if err != nil {
			return nil, err
		}
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
	mtp, err := i.state.Revocations.GenerateRevocationProof(rID, i.state.CommittedState.RevocationTreeRoot)
	if err != nil {
		return nil, err
	}
	res.MTP = mtp
	res.Issuer.RevocationTreeRoot = i.state.CommittedState.RevocationTreeRoot.Hex()
	res.Issuer.RootOfRoots = i.state.CommittedState.RootsTreeRoot.Hex()
	res.Issuer.ClaimsTreeRoot = i.state.CommittedState.ClaimsTreeRoot.Hex()

	stateHash, err := i.state.CommittedState.State()
	if err != nil {
		return nil, err
	}
	res.Issuer.State = stateHash.Hex()

	return res, nil
}

func (i *Identity) PublishLatestState(ctx context.Context) (string, error) {
	logger.Debug("PublishLatestState() invoked")

	publisher := Publisher{
		i:            i,
		circuitsPath: i.circuitsPath,
		stateStore:   i.stateStore,
	}
	inputs, err := publisher.PrepareInputs()
	if err != nil {
		return "", err
	}
	proof, err := publisher.GenerateProof(ctx, inputs)
	if err != nil {
		return "", err
	}

	latestState, err := i.state.CommittedState.State()
	if err != nil {
		return "", err
	}
	newState, err := i.state.GetStateHash()
	if err != nil {
		return "", err
	}

	if latestState.Equals(newState) {
		return "", errors.New("nothing to update")
	}

	ti := &TransitionInfoRequest{
		Identifier:        i.Identifier,
		LatestState:       latestState,
		NewState:          newState,
		IsOldStateGenesis: i.state.CommittedState.IsLatestStateGenesis,
		Proof:             proof.Proof,
	}
	txHex, err := publisher.UpdateState(ctx, ti)
	if err != nil {
		return "", err
	}
	logger.Info("transaction for change state:", txHex)

	return txHex, nil
}

func (i *Identity) sign(z *big.Int) ([]byte, error) {
	if !utils.CheckBigIntInField(z) {
		return nil, errors.New("data to signBytes is too large")
	}

	sig := i.sk.SignPoseidon(z).Compress()
	return sig[:], nil
}

func (i *Identity) getMTPProof(claimIdx *big.Int) (*verifiable.Iden3SparseMerkleProof, error) {
	mtpProof, _, err := i.state.Claims.Tree.GenerateProof(
		context.Background(), claimIdx, i.state.CommittedState.ClaimsTreeRoot)
	if err != nil {
		return nil, err
	}

	if i.state.CommittedState.Info.TxId == "" {
		return nil, errors.New("failed generate mtp proof. Transaction not exists")
	}

	txID := i.state.CommittedState.Info.TxId
	blockTimestamp := int(i.state.CommittedState.Info.BlockTimestamp)
	blockNumber := int(i.state.CommittedState.Info.BlockNumber)
	committedState, err := i.state.CommittedState.State()
	if err != nil {
		return nil, err
	}

	mtp := &verifiable.Iden3SparseMerkleProof{
		Type: verifiable.Iden3SparseMerkleProofType,
		IssuerData: verifiable.IssuerData{
			ID: i.Identifier,
			State: verifiable.State{
				TxID:               &txID,
				BlockTimestamp:     &blockTimestamp,
				BlockNumber:        &blockNumber,
				RootOfRoots:        strptr(i.state.CommittedState.RootsTreeRoot.Hex()),
				ClaimsTreeRoot:     strptr(i.state.CommittedState.ClaimsTreeRoot.Hex()),
				RevocationTreeRoot: strptr(i.state.CommittedState.RevocationTreeRoot.Hex()),
				Value:              strptr(committedState.Hex()),
			},
		},
		MTP: mtpProof,
	}

	return mtp, nil
}

func strptr(s string) *string {
	return &s
}
