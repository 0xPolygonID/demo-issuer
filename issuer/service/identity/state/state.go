package state

import (
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"issuer/db"
	"issuer/db/models"
	"issuer/service/claim"
)

type IdentityState struct {
	Claims      *Claims
	Revocations *Revocations
	Roots       *Roots
	db          *db.DB
}

func NewIdentityState(db *db.DB, treeDepth int) (*IdentityState, error) {
	logger.Debug("creating new identity state")

	claims, err := NewClaims(db, treeDepth)
	if err != nil {
		return nil, err
	}

	revs, err := NewRevocations(db, treeDepth)
	if err != nil {
		return nil, err
	}

	roots, err := NewRoots(db, treeDepth)
	if err != nil {
		return nil, err
	}

	return &IdentityState{
		Claims:      claims,
		Revocations: revs,
		Roots:       roots,
		db:          db,
	}, nil
}

func (is *IdentityState) SaveIdentity(identifier string) error {
	state, err := is.GetStateHash()
	if err != nil {
		return err
	}

	return is.db.SaveIdentity(&models.Identity{
		Identifier:         identifier,
		RootOfRoots:        is.Roots.Tree.Root().Hex(),
		RevocationTreeRoot: is.Revocations.Tree.Root().Hex(),
		ClaimsTreeRoot:     is.Claims.Tree.Root().Hex(),
		State:              state.Hex(),
	})
}

func (is *IdentityState) AddClaimToTree(c *core.Claim) error {
	logger.Debug("IdentityState.AddClaimToTree() invoked")

	return is.Claims.SaveClaimMT(c)
}

func (is *IdentityState) AddClaimToDB(c *claim.Claim) error {
	logger.Debug("IdentityState.AddClaimToDB() invoked")

	return is.Claims.SaveClaimDB(c)
}

func (is *IdentityState) GetStateHash() (*merkletree.Hash, error) {
	logger.Debug("GetStateHash() invoked")

	return merkletree.HashElems(
		is.Claims.Tree.Root().BigInt(),
		is.Revocations.Tree.Root().BigInt(),
		is.Roots.Tree.Root().BigInt(),
	)
}
