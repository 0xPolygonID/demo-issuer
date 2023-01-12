package state

import (
	"context"
	"math/big"

	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
)

type Revocations struct {
	Tree *merkletree.MerkleTree
}

func NewRevocations(treeStorage *store.BoltStore, treeDepth int) (*Revocations, error) {
	logger.Debug("creating new revocations state")

	revsTree, err := merkletree.NewMerkleTree(context.Background(), treeStorage.WithPrefix([]byte("revocation")), treeDepth)
	if err != nil {
		return nil, err
	}

	return &Revocations{
		Tree: revsTree,
	}, nil

}

func (r *Revocations) GenerateRevocationProof(nonce *big.Int, root *merkletree.Hash) (*merkletree.Proof, error) {
	logger.Debugf("GenerateRevocationProof() invoked with nonce of %d", nonce)

	proof, _, err := r.Tree.GenerateProof(context.Background(), nonce, root)
	return proof, err
}
