package contract

import "github.com/iden3/go-merkletree-sql"

// GetRevocationStatus status of revocation nonce. Info required to check revocation state of claim in circuits
type GetRevocationStatus struct {
	Issuer struct {
		State              *string `json:"state,omitempty"`
		RootOfRoots        *string `json:"root_of_roots,omitempty"`
		ClaimsTreeRoot     *string `json:"claims_tree_root,omitempty"`
		RevocationTreeRoot *string `json:"revocation_tree_root,omitempty"`
	} `json:"issuer"`
	MTP merkletree.Proof `json:"mtp"`
}
