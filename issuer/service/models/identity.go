package models

import core "github.com/iden3/go-iden3-core"

// Identity model
type Identity struct {
	Identifier string
	State      IdentityState
}

// NewIdentityFromIdentifier default identity model from identity and root state
func NewIdentityFromIdentifier(id *core.ID, rootState string) *Identity {
	return &Identity{
		Identifier: id.String(),
		State: IdentityState{
			Identifier: id.String(),
			State:      &rootState,
		},
	}
}
