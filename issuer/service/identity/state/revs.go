package state

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	logger "github.com/sirupsen/logrus"
	"math/big"
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

// GenerateRevocationProof generates the proof of existence (or non-existence) of an nonce in RevocationTree
func (r *Revocations) GenerateRevocationProof(nonce *big.Int) (*merkletree.Proof, error) {
	logger.Debugf("GenerateRevocationProof() invoked with nonce of %d", nonce)

	proof, _, err := r.Tree.GenerateProof(context.Background(), nonce, r.Tree.Root())
	return proof, err
}
