package models

import (
	"database/sql/driver"
	"strconv"

	"github.com/iden3/go-merkletree-sql"
)

// RevNonceUint64 type for revocation nonce. As go SQL driver do not support storing uint64 we have to use a trick with
// custom type to make it possible to store uint64 field in NUMERIC field in the PostgresDB
type RevNonceUint64 uint64

// Value implementation of valuer interface to convert uint64 value for storing in Postgres
func (r RevNonceUint64) Value() (driver.Value, error) {
	return strconv.FormatUint(uint64(r), 10), nil
}

// RevStatus status of revocation nonce
type RevStatus int

// Revocation model
type Revocation struct {
	ID          int64
	Identifier  string
	Nonce       RevNonceUint64
	Version     uint32
	Status      RevStatus
	Description string
}

// RevocationStatus status of revocation nonce. Info required to check revocation state of claim in circuits
type RevocationStatus struct {
	Issuer struct {
		State              *string `json:"state"`
		RootOfRoots        *string `json:"root_of_roots,omitempty"`
		ClaimsTreeRoot     *string `json:"claims_tree_root,omitempty"`
		RevocationTreeRoot *string `json:"revocation_tree_root,omitempty"`
	} `json:"issuer"`
	MTP *merkletree.Proof `json:"mtp"`
}
