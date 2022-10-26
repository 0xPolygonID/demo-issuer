package contract

import "github.com/iden3/go-merkletree-sql"

// GetRevocationStatusResponse status of revocation nonce. Info required to check revocation state of claim in circuits
type GetRevocationStatusResponse struct {
	Issuer struct {
		State              string `codec:"state,omitempty"`
		RootOfRoots        string `codec:"root_of_roots,omitempty"`
		ClaimsTreeRoot     string `codec:"claims_tree_root,omitempty"`
		RevocationTreeRoot string `codec:"revocation_tree_root,omitempty"`
	} `codec:"issuer"`
	MTP *merkletree.Proof `codec:"mtp"`
}
