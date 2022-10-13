package models

// Schema base model for claim
type Schema struct {
	ID      int64
	Name    string
	Encoded string
	Format  SchemaFormat
	URL     string
}

// SchemaFormat type
type SchemaFormat string

const (
	// JSONLD JSON-LD schema format
	JSONLD SchemaFormat = "json-ld"

	// JSON JSON schema format
	JSON SchemaFormat = "json"
)

const (
	AuthBJJCredentialHash = "ca938857241db9451ea329256b9c06e5"
	AuthBJJCredential     = "AuthBJJCredential"
	AuthBJJCredentialURL  = "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/auth.json-ld"
)
