package contract

import "encoding/json"

type CreateClaimRequest struct {
	Schema          *Schema         `codec:"schema"`
	Data            json.RawMessage `codec:"data"`
	Identifier      string          `codec:"identifier"`
	Expiration      int64           `codec:"expiration"`
	Version         uint32          `codec:"version"`
	RevNonce        *uint64         `codec:"revNonce"`
	SubjectPosition string          `codec:"subjectPosition"`
}

type Schema struct {
	URL  string `codec:"url"`
	Type string `codec:"type"`
}
