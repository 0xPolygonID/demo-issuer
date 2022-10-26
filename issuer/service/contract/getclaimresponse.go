package contract

import "time"

type GetClaimResponse struct {
	ID                string                 `json:"id"`
	Context           []string               `json:"@context"`
	Type              []string               `json:"@type"`
	Expiration        time.Time              `json:"expiration,omitempty"`
	Updatable         bool                   `json:"updatable"`
	Version           uint32                 `json:"version"`
	RevNonce          uint64                 `json:"rev_nonce"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	CredentialStatus  *CredentialStatus      `json:"credentialStatus,omitempty"`
	SubjectPosition   string                 `json:"subject_position,omitempty"`
	CredentialSchema  *CredentialStatus      `json:"credentialSchema"`
	Proof             interface{}            `json:"proof,omitempty"`
}

type CredentialStatus struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type CredentialStatusType string
