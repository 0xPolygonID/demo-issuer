package identitystate

import (
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
	"issuer/service/claims"
	"issuer/service/revocations"
	"issuer/service/roots"
)

type IdentityState struct {
	Claims      *claims.Claims
	Revocations *revocations.Revocations
	Roots       *roots.Roots
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
	}, nil
}

func (is *IdentityState) GetStateHash() (*merkletree.Hash, error) {
	return merkletree.HashElems(
		is.Claims.ClaimTree.Root().BigInt(),
		is.Revocations.RevTree.Root().BigInt(),
		is.Roots.RootsTree.Root().BigInt(),
	)
}
