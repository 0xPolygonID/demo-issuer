package models

import (
	"github.com/iden3/go-circuits"
)

// IdentityStatus represents type for state Status
type IdentityStatus string

// IdentityState identity state model
type IdentityState struct {
	StateID            int64
	Identifier         string
	State              *string
	RootOfRoots        *string
	ClaimsTreeRoot     *string
	RevocationTreeRoot *string
}

// ToTreeState returns circuits.TreeState structure
func (is IdentityState) ToTreeState() (circuits.TreeState, error) {
	return BuildTreeState(is.State, is.ClaimsTreeRoot, is.RevocationTreeRoot,
		is.RootOfRoots)
}
