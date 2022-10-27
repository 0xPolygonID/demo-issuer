package roots

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
)

type Roots struct {
	RootsTree *merkletree.MerkleTree
}

func New(db *db.DB, treeDepth int) (*Roots, error) {

	treeStorage, err := store.NewBoltStorage(db.GetConnection())
	if err != nil {
		return nil, err
	}

	roots, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("ror")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Roots{
		RootsTree: roots,
	}, nil
}