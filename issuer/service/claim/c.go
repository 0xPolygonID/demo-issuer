package claim

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/iden3/go-schema-processor/verifiable"
	"issuer/service/models"
	"issuer/service/schema"
	"issuer/service/utils"
	"strings"
	"time"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"
)

type CoreClaimData struct {
	EncodedSchema   string
	Slots           processor.ParsedSlots
	SubjectID       string
	Expiration      int64
	Version         uint32
	Nonce           *uint64
	SubjectPosition string
}

// GenerateCoreClaim generate core claim via settings from CoreClaimData.
func GenerateCoreClaim(req CoreClaimData) (*core.Claim, error) {
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

func CoreClaimToClaimModel(claim *core.Claim, schemaURL, schemaType string) (*models.Claim, error) {
	otherIdentifier := ""
	id, err := claim.GetID()
	if err != nil {
		return nil, err
	}
	otherIdentifier = id.String()

	var expiration int64
	if expirationDate, ok := claim.GetExpirationDate(); ok {
		expiration = expirationDate.Unix()
	}

	hindex, err := claim.HIndex()
	if err != nil {
		return nil, err
	}

	sb := claim.GetSchemaHash()
	schemaHash := hex.EncodeToString(sb[:])
	res := models.Claim{
		SchemaHash:      schemaHash,
		SchemaURL:       schemaURL,
		SchemaType:      schemaType,
		OtherIdentifier: otherIdentifier,
		Expiration:      expiration,
		Updatable:       claim.GetFlagUpdatable(),
		Version:         claim.GetVersion(),
		RevNonce:        claim.GetRevocationNonce(),
		CoreClaim:       claim,
		HIndex:          hindex.String(),
	}

	return &res, nil
}

func NewAuthClaim(key *babyjub.PublicKey, schemaHash core.SchemaHash) (*core.Claim, error) {
	revNonce, err := utils.Rand()
	if err != nil {
		return nil, err
	}

	return core.NewClaim(schemaHash,
		core.WithIndexDataInts(key.X, key.Y),
		core.WithRevocationNonce(revNonce))
}

func CreateCredentialStatus(urlBase string, sType verifiable.CredentialStatusType, revNonce uint64) ([]byte, error) {
	cStatus := verifiable.CredentialStatus{
		ID:   fmt.Sprintf("%s/api/v1/claims/revocation/status/%d", urlBase, revNonce),
		Type: sType,
	}

	b, err := json.Marshal(cStatus)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ClaimModelToIden3Credential(c *models.Claim) (*verifiable.Iden3Credential, error) {
	claimIdPos, err := getClaimIdPosition(c.CoreClaim)
	if err != nil {
		return nil, err
	}

	credSubjects := make(map[string]interface{})
	err = json.Unmarshal(c.Data, &credSubjects)
	if err != nil {
		return nil, err
	}
	credSubjects["type"] = c.SchemaType
	if len(c.OtherIdentifier) > 0 {
		credSubjects["id"] = c.OtherIdentifier
	}

	// * create proof object
	proofs := make([]interface{}, 0)

	signatureProof := &verifiable.BJJSignatureProof2021{}
	if c.SignatureProof != nil && string(c.CredentialStatus) != "{}" {
		err = json.Unmarshal(c.SignatureProof, signatureProof)
		if err != nil {
			return nil, err
		}

		proofs = append(proofs, signatureProof)
	}

	mtpProof := &verifiable.Iden3SparseMerkleProof{}
	if !strings.EqualFold(string(c.MTPProof), "{}") {
		err = json.Unmarshal(c.MTPProof, mtpProof)
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, mtpProof)
	}

	// create credential status object
	credStatus := &verifiable.CredentialStatus{}
	if c.CredentialStatus != nil && string(c.CredentialStatus) != "{}" {
		err = json.Unmarshal(c.CredentialStatus, credStatus)
		if err != nil {
			return nil, err
		}
	}

	res := &verifiable.Iden3Credential{
		ID:         string(c.ID),
		Type:       []string{schema.Iden3CredentialSchema},
		Expiration: time.Unix(c.Expiration, 0),
		RevNonce:   c.RevNonce,
		Updatable:  c.Updatable,
		Version:    c.Version,
		Context:    []string{schema.Iden3CredentialSchemaURL, c.SchemaURL},
		CredentialSchema: struct {
			ID   string `json:"@id"`
			Type string `json:"type"`
		}{
			ID:   c.SchemaURL,
			Type: c.SchemaType,
		},
		SubjectPosition:   claimIdPos,
		CredentialSubject: credSubjects,
		Proof:             proofs,
		CredentialStatus:  credStatus,
	}

	return res, nil
}

func getClaimIdPosition(c *core.Claim) (string, error) {
	claimIdPos, err := c.GetIDPosition()
	if err != nil {
		return "", err
	}

	return subjectPositionIDToString(claimIdPos)
}

func subjectPositionIDToString(p core.IDPosition) (string, error) {
	switch p {
	case core.IDPositionIndex:
		return SubjectPositionIndex, nil
	case core.IDPositionValue:
		return SubjectPositionValue, nil
	default:
		return "", fmt.Errorf("id position is not specified")
	}
}
