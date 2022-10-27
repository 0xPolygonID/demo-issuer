package models

import (
	core "github.com/iden3/go-iden3-core"
)

type Claim struct {
	ID              []byte
	Identifier      *string
	Issuer          string
	SchemaHash      string
	SchemaURL       string
	SchemaType      string
	OtherIdentifier string
	Expiration      int64
	// TODO(illia-korotia): delete from db but left in struct.
	Updatable        bool
	Revoked          bool
	Version          uint32
	RevNonce         uint64
	Data             []byte
	CoreClaim        *core.Claim
	MTPProof         []byte
	SignatureProof   []byte
	IdentityState    *string
	Status           *IdentityStatus
	CredentialStatus []byte
	HIndex           string
}
