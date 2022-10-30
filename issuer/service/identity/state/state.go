package state

import (
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
	"issuer/service/claim"
	"issuer/service/models"
)

type IdentityState struct {
	Claims      *Claims
	Revocations *Revocations
	Roots       *Roots
	db          *db.DB
}

func NewIdentityState(db *db.DB, treeDepth int) (*IdentityState, error) {
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
	return is.Claims.SaveClaimMT(c)
}

func (is *IdentityState) AddClaimToDB(c *claim.Claim) error {
	return is.Claims.SaveClaimDB(c)
}

func (is *IdentityState) GetStateHash() (*merkletree.Hash, error) {
	return merkletree.HashElems(
		is.Claims.Tree.Root().BigInt(),
		is.Revocations.Tree.Root().BigInt(),
		is.Roots.Tree.Root().BigInt(),
	)
}
