package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-schema-processor/processor"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/pkg/errors"
	"lightissuer/app/configs"
	"lightissuer/app/rest"
	"lightissuer/models"
	"lightissuer/services"
	"lightissuer/utils"
	"net/http"
	"strconv"
	"time"
)

const (
	// SubjectPositionIndex save subject in index part of claim. By default.
	SubjectPositionIndex = "index"
	// SubjectPositionValue save subject in value part of claim.
	SubjectPositionValue = "value"
)

// ClaimRequest contains settings for issuing a claim.
type ClaimRequest struct {
	EncodedSchema   string
	Slots           processor.ParsedSlots
	SubjectID       string
	Expiration      int64
	Version         uint32
	Nonce           *uint64
	SubjectPosition string
}

// ClaimHandler provide handlers for claim requests
type ClaimHandler struct {
	dbService       *services.DBService
	service         *services.ClaimService
	schService      *services.SchemaService
	identityService *services.IdentityService
	config          *configs.Environement
}

// NewClaimHandler creates ClaimHandler
func NewClaimHandler(storage *services.DBService, service *services.ClaimService, schemaService *services.SchemaService, identityService *services.IdentityService, config *configs.Environement) *ClaimHandler {
	return &ClaimHandler{
		storage,
		service,
		schemaService,
		identityService,
		config,
	}
}

type claimReq struct {
	Schema struct {
		URL  string `json:"url"`
		Type string `json:"type"`
	} `json:"schema"`
	Data            json.RawMessage `json:"data"`
	Identifier      string          `json:"identifier"`
	Expiration      int64           `json:"expiration"`
	Version         uint32          `json:"version"`
	RevNonce        *uint64         `json:"revNonce"`
	SubjectPosition string          `json:"subjectPosition"`
}

// CreateClaim POST /api/v1/identities/{identifier}/claims
func (c *ClaimHandler) CreateClaim(w http.ResponseWriter, r *http.Request) {

	issuerID, err := core.IDFromString(c.identityService.Identity.Identifier)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't generate identity from the requested identity", 0)
	}

	var cReq claimReq
	if err := render.DecodeJSON(r.Body, &cReq); err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't bind request", 0)
		return
	}

	schemaBytes, _, err := c.schService.Load(r.Context(), cReq.Schema.URL)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, errors.Errorf("can't load schema: %s ", cReq.Schema.URL), "", 0)
		return
	}

	slots, err := c.schService.Process(r.Context(), cReq.Schema.URL, cReq.Schema.Type, cReq.Data)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, fmt.Sprintf("can't process claim data: %s", cReq.Schema), 0)
		return
	}

	claimReq := ClaimRequest{
		EncodedSchema:   CreateSchemaHash(schemaBytes, cReq.Schema.Type),
		Slots:           slots,
		SubjectID:       cReq.Identifier,
		Expiration:      cReq.Expiration,
		Version:         cReq.Version,
		Nonce:           cReq.RevNonce,
		SubjectPosition: cReq.SubjectPosition,
	}

	coreClaim, err := GenerateCoreClaim(claimReq)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err,
			"can't create core claim", 0)
		return
	}

	claim, err := models.FromClaimer(coreClaim, cReq.Schema.URL, cReq.Schema.Type)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "can't transform claimer", 0)
		return
	}

	// set credential status
	issuerIDString := issuerID.String()
	err = claim.SetCredentialStatus(c.config.Host, issuerIDString, verifiable.SparseMerkleTreeProof)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "can't set credential status", 0)
		return
	}

	authClaim, err := c.service.GetAuthClaim(&issuerID)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "can't get auth claim of issuer", 0)
		return
	}

	proof, err := c.identityService.SignClaimEntry(authClaim, coreClaim)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "can't sign claim", 0)
		return
	}

	proof.IssuerData.RevocationStatus = fmt.Sprintf("%s/api/v1/claims/revocation/status/%d",
		c.config.Host, claim.RevNonce)

	// Save
	claim.Identifier = &issuerIDString
	claim.Issuer = issuerIDString

	jsonSignatureProof, err := json.Marshal(proof)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err,
			"can't save claim", 0)
		return
	}
	claim.SignatureProof = jsonSignatureProof
	claim.Data = string(cReq.Data)

	var claimModel *models.Claim
	claimModel, err = c.service.Save(claim)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError,
			rest.New5XXHttpErr(err, "can't save claim"), "", 0)
		return
	}

	render.Status(r, http.StatusAccepted)
	render.JSON(w, r, map[string]interface{}{"id": claimModel.ID})
}

//// GetClaimByID http://localhost:8001/api/v1/identities/{identifier}/claims/{id}
func (c *ClaimHandler) GetClaimByID(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		rest.ErrorJSON(w, r, http.StatusBadRequest, errors.New("no claim identifier"), "can't  parse claim id param", 0)
		return
	}

	claimID, err := uuid.Parse(idParam)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't  parse claim id param", 0)
		return
	}

	claimIDByte := []byte(claimID.String())
	claim, err := c.service.GetByID(claimIDByte)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusNotFound, err, fmt.Sprintf("can't get a claim %s", claimID.String()), 0)
		return
	}

	cred, err := c.schService.FromClaimToIden3Credential(*claim)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, fmt.Sprintf("can't create JSON-LD credential: %d", claim.ID), 0)
		return
	}

	fmt.Printf("cred %+v", cred)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, cred)
}

func (c *ClaimHandler) RevocationStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Inside Revocation Handler")

	issuerID, err := core.IDFromString(c.identityService.Identity.Identifier)

	nonce, err := strconv.ParseUint(chi.URLParam(r, "nonce"), 10, 64)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't parse nonce", 0)
		return
	}

	ctx := r.Context()
	proof, err := c.service.GetRevocationNonceMTP(ctx, nonce)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, err, fmt.Sprintf("can't generate non revocation proof for revocation nonce: %d of issuer %s", nonce, issuerID.String()), 0)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, proof)
}

// GenerateCoreClaim generate core claim via settings from ClaimRequest.
func GenerateCoreClaim(req ClaimRequest) (*core.Claim, error) {
	var revNonce *uint64

	r, err := utils.RandInt64()
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
			return nil, errors.New("unknown subject position")
		}
	}

	if req.Expiration != 0 {
		coreClaim.SetExpirationDate(time.Unix(req.Expiration, 0))
	}
	return coreClaim, nil
}

// CreateSchemaHash returns hash of shema and credential type parameters
func CreateSchemaHash(schemaBytes []byte, credentialType string) string {
	var sHash core.SchemaHash
	h := crypto.Keccak256(schemaBytes, []byte(credentialType))
	copy(sHash[:], h[len(h)-16:])
	return hex.EncodeToString(sHash[:])
}
