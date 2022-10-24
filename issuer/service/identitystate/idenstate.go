package identitystate

import (
	"github.com/iden3/go-merkletree-sql"
	"issuer/service/claims"
	"issuer/service/contract"
	"issuer/service/db"
	"issuer/service/models"
	"issuer/service/revocations"
	"issuer/service/roots"
	"math/big"
)

type IdentityState struct {
	Claims      *claims.Claims
	Revocations *revocations.Revocations
	Roots       *roots.Roots

	//State string
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

func (is *IdentityState) GetIdentityState() (*contract.IdentityState, error) {

}

func (is *IdentityState) GetClaim() (*contract.IdentityState, error) {

}

func (is *IdentityState) GetRevocationStatus(nonce uint64) (*models.RevocationStatus, error) {
	rID := new(big.Int).SetUint64(nonce)

	res := &models.RevocationStatus{}
	mtp, err := is.Revocations.GenerateRevocationProof(rID)
	if err != nil {
		return nil, err
	}
	res.MTP = mtp
	res.Issuer.RevocationTreeRoot = is.Revocations.State
	res.Issuer.RootOfRoots = is.Roots.RootsTree
	res.Issuer.ClaimsTreeRoot = is.Claims.ClaimTree.Root().Hex()

	stateHash, err := is.GetStateHash()
	if err != nil {
		return nil, err
	}
	res.Issuer.State = stateHash.Hex()

	return res, nil
}

func (is *IdentityState) SaveClaim(claim *models.Claim) error {

	// save to the tree (Add later)

	// save to the db
	return is.Claims.SaveClaim(claim)
}
