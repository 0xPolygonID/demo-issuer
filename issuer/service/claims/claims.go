package claims

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	"issuer/service/db"
	"issuer/service/models"
)

type Claims struct {
	db        *db.DB
	ClaimTree *merkletree.MerkleTree
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
		db:        db,
		ClaimTree: claimTree,
	}, nil
}

func (c *Claims) GetClaim(id []byte) (*models.Claim, error) {
	cl, err := c.db.GetClaim(id)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

//func (c *Claims) GetClaim(id string) (*models.Claim, error) {
//	idBytes := []byte(id)
//	cl, err := c.db.GetClaim(idBytes)
//	if err != nil {
//		return nil, err
//	}
//
//	return cl, nil
//}

func (c *Claims) SaveClaim(claim *models.Claim) error {
	return c.db.SaveClaim(claim)
}
