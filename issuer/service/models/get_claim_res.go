package models

import (
	"github.com/iden3/go-schema-processor/verifiable"
)

type GetClaimResponse *verifiable.Iden3Credential

//type Iden3Credential struct {
//	ID                string                 `json:"id"`
//	Context           []string               `json:"@context"`
//	Type              []string               `json:"@type"`
//	Expiration        time.Time              `json:"expiration,omitempty"`
//	Updatable         bool                   `json:"updatable"`
//	Version           uint32                 `json:"version"`
//	RevNonce          uint64                 `json:"rev_nonce"`
//	CredentialSubject map[string]interface{} `json:"credentialSubject"`
//	CredentialStatus  *CredentialStatus      `json:"credentialStatus,omitempty"`
//	SubjectPosition   string                 `json:"subject_position,omitempty"`
//	CredentialSchema  struct {
//		ID   string `json:"@id"`
//		Type string `json:"type"`
//	} `json:"credentialSchema"`
//	Proof interface{} `json:"proof,omitempty"`
//}
