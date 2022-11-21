package state

import (
	store "github.com/demonsh/smt-bolt"
	"github.com/google/uuid"
	"github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"issuer/db"
	"issuer/service/claim"
	"issuer/service/schema"
)

type IdentityState struct {
	IsGenesis bool
	State     merkletree.Hash
	ClR       merkletree.Hash
	ReR       merkletree.Hash
	RoR       merkletree.Hash

	ClaimsTree      *Claims
	RevocationsTree *Revocations
	RootsTree       *Roots
	db              *db.DB
}

func NewIdentityState(db *db.DB, treeDepth int) (*IdentityState, error) {
	logger.Debug("creating new identity state")

	treeStorage, err := store.NewBoltStorage(db.GetConnection())
	if err != nil {
		return nil, err
	}

	claims, err := NewClaims(db, treeStorage, treeDepth)
	if err != nil {
		return nil, err
	}

	revs, err := NewRevocations(treeStorage, treeDepth)
	if err != nil {
		return nil, err
	}

	roots, err := NewRoots(treeStorage, treeDepth)
	if err != nil {
		return nil, err
	}

	idenState := &IdentityState{
		IsGenesis:       true,
		ClaimsTree:      claims,
		RevocationsTree: revs,
		RootsTree:       roots,
		db:              db,
	}
	err = idenState.UpdateState()
	if err != nil {
		return nil, err
	}
	return idenState, nil
}

func (is *IdentityState) ToCircuitTreeState() (circuits.TreeState, error) {
	return circuits.TreeState{
		State:          &is.State,
		ClaimsRoot:     &is.ClR,
		RevocationRoot: &is.ReR,
		RootOfRoots:    &is.RoR,
	}, nil
}

func (is *IdentityState) SetupGenesisState(pk *babyjub.PublicKey) (*core.ID, *core.Claim, error) {
	logger.Trace("getting auth schema hash")

	schemaHash, err := core.NewSchemaHashFromHex(schema.AuthBJJCredentialHash)
	if err != nil {
		return nil, nil, err
	}

	logger.Trace("creating new auth claim")
	authClaim, err := claim.NewAuthClaim(pk, schemaHash)
	if err != nil {
		return nil, nil, err
	}

	logger.Trace("adding auth claim to the claims tree")
	err = is.AddClaimToTree(authClaim)
	if err != nil {
		return nil, nil, err
	}

	currState := is.State
	if err != nil {
		return nil, nil, err
	}

	identifier, err := core.IdGenesisFromIdenState(core.TypeDefault, currState.BigInt())
	if err != nil {
		return nil, nil, err
	}

	return identifier, authClaim, nil
}

func (is *IdentityState) SaveIdentity(identifier *core.ID, authClaimId uuid.UUID) error {

	id := identifier.Bytes()
	authId := []byte(authClaimId.String())

	return is.db.SaveIdentity(id, authId)

}

func (is *IdentityState) GetIdentityFromDB() (*core.ID, *uuid.UUID, error) {
	logger.Debug("IdentityState.GetIdentityFromDB() invoked")

	id, authClaimId, err := is.db.GetSavedIdentity()
	if err != nil {
		return nil, nil, err
	}

	if id == nil {
		return nil, nil, nil
	}
	coreId, err := core.IDFromBytes(id)
	if err != nil {
		return nil, nil, err
	}

	claimId, err := uuid.Parse(string(authClaimId))
	if err != nil {
		return nil, nil, err
	}

	return &coreId, &claimId, nil
}

func (is *IdentityState) AddClaimToTree(c *core.Claim) error {
	logger.Debug("IdentityState.AddClaimToTree() invoked")

	return is.ClaimsTree.SaveClaimMT(c)
}

func (is *IdentityState) AddClaimToDB(c *claim.Claim) error {
	logger.Debug("IdentityState.AddClaimToDB() invoked")

	return is.ClaimsTree.SaveClaimDB(c)
}

func (is *IdentityState) UpdateState() error {
	logger.Debug("IdentityState.UpdateState() invoked")

	is.ClR = *is.ClaimsTree.Tree.Root()
	is.ReR = *is.RevocationsTree.Tree.Root()
	is.RoR = *is.RootsTree.Tree.Root()

	ns, err := merkletree.HashElems(is.ClR.BigInt(), is.ReR.BigInt(), is.RoR.BigInt())
	if err != nil {
		return err
	}
	is.State = *ns
	return nil
}
