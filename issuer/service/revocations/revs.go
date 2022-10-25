package revocations

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	"issuer/db"
	"math/big"
)

type Revocations struct {
	RevTree *merkletree.MerkleTree
}

func New(db *db.DB, treeDepth int) (*Revocations, error) {
	treeStorage, err := store.NewBoltStorage(db.GetConnection())
	if err != nil {
		return nil, err
	}

	revsTree, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("revocation")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Revocations{
		RevTree: revsTree,
	}, nil

}

// GenerateRevocationProof generates the proof of existence (or non-existence) of an nonce in RevocationTree
func (r *Revocations) GenerateRevocationProof(nonce *big.Int) (*merkletree.Proof, error) {
	proof, _, err := r.RevTree.GenerateProof(context.Background(), nonce, r.RevTree.Root())
	return proof, err
}
