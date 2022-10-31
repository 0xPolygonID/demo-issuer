package state

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"issuer/db"
)

type Roots struct {
	Tree *merkletree.MerkleTree
}

func NewRoots(db *db.DB, treeDepth int) (*Roots, error) {
	logger.Debug("creating new roots state")

	treeStorage, err := store.NewBoltStorage(db.GetConnection())
	if err != nil {
		return nil, err
	}

	roots, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("ror")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Roots{
		Tree: roots,
	}, nil
}
