package models

import (
	"github.com/iden3/go-circuits"
	"time"
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
	BlockTimestamp     *int
	BlockNumber        *int
	TxID               *string
	PreviousState      *string
	Status             *IdentityStatus
	ModifiedAt         time.Time
	CreatedAt          time.Time
}

// ToTreeState returns circuits.TreeState structure
func (is IdentityState) ToTreeState() (circuits.TreeState, error) {
	return BuildTreeState(is.State, is.ClaimsTreeRoot, is.RevocationTreeRoot,
		is.RootOfRoots)
}
