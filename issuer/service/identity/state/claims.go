package state

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"issuer/service/claim"
	"issuer/service/db"
)

type Claims struct {
	db   *db.DB
	Tree *merkletree.MerkleTree
}

func NewClaims(db *db.DB, treeStorage *store.BoltStore, treeDepth int) (*Claims, error) {
	logger.Debug("creating new claims state")

	claimTree, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("claims")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Claims{
		db:   db,
		Tree: claimTree,
	}, nil
}

func (c *Claims) GetClaim(id []byte) (*claim.Claim, error) {
	logger.Debugf("GetClaim() invoked with id %x", id)

	cl, err := c.db.GetClaim(id)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *Claims) SaveClaimDB(claim *claim.Claim) error {
	logger.Debugf("SaveClaimDB() invoked with claim %v", claim)

	return c.db.SaveClaim(claim)
}

func (c *Claims) SaveClaimMT(claim *core.Claim) error {
	logger.Debugf("SaveClaimMT() invoked with claim %v", claim)

	i, v, err := claim.HiHv()
	if err != nil {
		return err
	}

	return c.Tree.Add(context.Background(), i, v)
}
