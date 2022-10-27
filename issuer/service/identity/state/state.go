package state

import (
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
	"issuer/service/identity/state/claims"
	"issuer/service/identity/state/revocations"
	"issuer/service/identity/state/roots"
	"issuer/service/models"
)

type IdentityState struct {
	Claims      *claims.Claims
	Revocations *revocations.Revocations
	Roots       *roots.Roots
	db          *db.DB
}

func New(db *db.DB, treeDepth int) (*IdentityState, error) {
	claims, err := claims.New(db, treeDepth)
	if err != nil {
		return nil, err
	}

	revs, err := revocations.New(db, treeDepth)
	if err != nil {
		return nil, err
	}

	roots, err := roots.New(db, treeDepth)
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

func (is *IdentityState) GetStateHash() (*merkletree.Hash, error) {
	return merkletree.HashElems(
		is.Claims.Tree.Root().BigInt(),
		is.Revocations.Tree.Root().BigInt(),
		is.Roots.Tree.Root().BigInt(),
	)
}
