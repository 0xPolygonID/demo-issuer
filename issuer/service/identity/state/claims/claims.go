package claims

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
	"issuer/service/models"
)

type Claims struct {
	db   *db.DB
	Tree *merkletree.MerkleTree
}

func New(db *db.DB, treeDepth int) (*Claims, error) {
	treeStorage, err := store.NewBoltStorage(db.GetConnection())
	if err != nil {
		return nil, err
	}

	claimTree, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("claims")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Claims{
		db:   db,
		Tree: claimTree,
	}, nil
}

func (c *Claims) GetClaim(id []byte) (*models.Claim, error) {
	cl, err := c.db.GetClaim(id)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (c *Claims) SaveClaimDB(claim *models.Claim) error {
	return c.db.SaveClaim(claim)
}

func (c *Claims) SaveClaimMT(claim *core.Claim) error {
	i, v, err := claim.HiHv()
	if err != nil {
		return err
	}

	return c.Tree.Add(context.Background(), i, v)
}
