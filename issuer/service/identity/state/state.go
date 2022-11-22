package state

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"issuer/db"
	"issuer/service/claim"
	"issuer/service/schema"
	"math/big"
)

// Info contains information about when the state was committed.
type Info struct {
	TxId           string
	BlockTimestamp uint64
	BlockNumber    uint64
}

type CommittedState struct {
	Info *Info

	IsLatestStateGenesis bool
	RootsTreeRoot        *merkletree.Hash
	ClaimsTreeRoot       *merkletree.Hash
	RevocationTreeRoot   *merkletree.Hash
}

func (cs *CommittedState) State() (*merkletree.Hash, error) {
	return merkletree.HashElems(cs.ClaimsTreeRoot.BigInt(), cs.RevocationTreeRoot.BigInt(), cs.RootsTreeRoot.BigInt())
}

type IdentityState struct {
	CommittedState CommittedState

	Claims      *Claims
	Revocations *Revocations
	Roots       *Roots
	db          *db.DB
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

	return &IdentityState{
		Claims:      claims,
		Revocations: revs,
		Roots:       roots,
		db:          db,
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

	currState, err := is.GetStateHash()
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

func (is *IdentityState) IsGenesis() bool {
	return is.Claims.Tree.Root().Equals(&merkletree.HashZero) &&
		is.Roots.Tree.Root().Equals(&merkletree.HashZero) &&
		is.Revocations.Tree.Root().Equals(&merkletree.HashZero)
}

func (is *IdentityState) GetInclusionProof(claim *core.Claim) (*merkletree.Proof, *big.Int, error) {
	hi, _, err := claim.HiHv()
	if err != nil {
		return nil, nil, err
	}
	return is.Claims.Tree.GenerateProof(context.Background(), hi, is.CommittedState.ClaimsTreeRoot)
}

func (is *IdentityState) GetRevocationProof(claim *core.Claim) (*merkletree.Proof, *big.Int, error) {
	hi, _, err := claim.HiHv()
	if err != nil {
		return nil, nil, err
	}
	return is.Revocations.Tree.GenerateProof(context.Background(), hi, is.CommittedState.RevocationTreeRoot)
}
