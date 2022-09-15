package services

import (
	"context"
	store "github.com/demonsh/smt-bolt"
	"github.com/iden3/go-merkletree-sql"
	"github.com/pkg/errors"
	"math/big"
)

// IdentityMerkleTrees structure stores loaded merkle trees for one Identity
type IdentityMerkleTrees struct {
	ClaimTree *merkletree.MerkleTree
	RevTree   *merkletree.MerkleTree
	RorTree   *merkletree.MerkleTree
}

const mtDepth = 32

func CreateIdentityMerkleTrees(ctx context.Context, treeStorage *store.BoltStore) (*IdentityMerkleTrees, error) {

	claimStorage := treeStorage.WithPrefix([]byte("claims"))
	revStorage := treeStorage.WithPrefix([]byte("revocation"))
	rorStorage := treeStorage.WithPrefix([]byte("ror"))

	claimTree, err := merkletree.NewMerkleTree(ctx, claimStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	revTree, err := merkletree.NewMerkleTree(ctx, revStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	rorTree, err := merkletree.NewMerkleTree(ctx, rorStorage, mtDepth)
	if err != nil {
		return nil, err
	}

	imts := &IdentityMerkleTrees{
		claimTree,
		revTree,
		rorTree,
	}
	return imts, nil
}

// GenerateRevocationProof generates the proof of existence (or non-existence) of an nonce in RevocationTree
func (imts *IdentityMerkleTrees) GenerateRevocationProof(ctx context.Context,
	nonce *big.Int, root *merkletree.Hash) (*merkletree.Proof, error) {
	proof, _, err := imts.RevTree.GenerateProof(ctx, nonce, root)
	return proof, errors.WithStack(err)
}
