package identity

import (
	"encoding/hex"
	"fmt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/processor"
	"issuer/service/utils"
	"time"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"
)

type ClaimRequest struct {
	EncodedSchema   string
	Slots           processor.ParsedSlots
	SubjectID       string
	Expiration      int64
	Version         uint32
	Nonce           *uint64
	SubjectPosition string
}

// GenerateCoreClaim generate core claim via settings from ClaimRequest.
func GenerateCoreClaim(req ClaimRequest) (*core.Claim, error) {
	var revNonce *uint64
	r, err := utils.Rand()
	if err != nil {
		return nil, err
	}

	revNonce = &r
	if req.Nonce != nil {
		revNonce = req.Nonce
	}

	var coreClaim *core.Claim

	var sh core.SchemaHash
	schemaBytes, err := hex.DecodeString(req.EncodedSchema)
	if err != nil {
		return nil, err
	}
	copy(sh[:], schemaBytes)
	coreClaim, err = core.NewClaim(sh,
		core.WithIndexDataBytes(req.Slots.IndexA, req.Slots.IndexB),
		core.WithValueDataBytes(req.Slots.ValueA, req.Slots.ValueB),
		core.WithRevocationNonce(*revNonce),
		core.WithVersion(req.Version))
	if err != nil {
		return nil, err
	}

	if req.SubjectID != "" {
		var userID core.ID
		userID, err = core.IDFromString(req.SubjectID)
		if err != nil {
			return nil, err
		}

		switch req.SubjectPosition {
		case "", SubjectPositionIndex:
			coreClaim.SetIndexID(userID)
		case SubjectPositionValue:
			coreClaim.SetValueID(userID)
		default:
			return nil, fmt.Errorf("unknown subject position")
		}
	}

	if req.Expiration != 0 {
		coreClaim.SetExpirationDate(time.Unix(req.Expiration, 0))
	}
	return coreClaim, nil
}
