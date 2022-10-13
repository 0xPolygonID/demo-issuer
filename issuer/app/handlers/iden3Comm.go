package handlers

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-jwz"
	"github.com/iden3/go-schema-processor/verifiable"
	"github.com/iden3/iden3comm"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"io"
	"lightissuer/app/rest"
	"lightissuer/models"
	"lightissuer/services"
	"net/http"
	"os"
	"path/filepath"
)

// MaxBodySizeBytes defines is 2 MB for protocol message
const MaxBodySizeBytes = 2 * 1000 * 1000

// IDEN3CommHandler provides handlers for zk proof
type IDEN3CommHandler struct {
	iService      *services.IdentityService
	dbService     *services.DBService
	claimService  *services.ClaimService
	schemaService *services.SchemaService
	circuitPath   string
}

// NewIDEN3CommHandler inits IDEN3Comm handler
func NewIDEN3CommHandler(
	iService *services.IdentityService,
	dbService *services.DBService,
	claimService *services.ClaimService,
	schemaService *services.SchemaService,
	circuitPath string,
) *IDEN3CommHandler {
	return &IDEN3CommHandler{
		iService:      iService,
		dbService:     dbService,
		claimService:  claimService,
		schemaService: schemaService,
		circuitPath:   circuitPath,
	}
}

// Handle POST /api/v1/agent
func (comm *IDEN3CommHandler) Handle(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodySizeBytes)
	envelope, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: define protocol errors
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't bind request to protocol message", 0)
		return
	}

	basicMessage, err := Unpack(comm.circuitPath, envelope)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "failed unpack protocol message", 0)
		return
	}

	// TODO(illia-korotia): validator like https://github.com/go-playground/validator?
	if basicMessage.To == "" {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "failed request. empty to field", 0)
		return
	}

	if basicMessage.From == "" {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "failed request. empty from field", 0)
		return
	}

	switch basicMessage.Type {
	case protocol.CredentialFetchRequestMessageType:
		fetchRequestBody := &protocol.CredentialFetchRequestMessageBody{}
		err = json.Unmarshal(basicMessage.Body, fetchRequestBody)
		if err != nil {
			rest.ErrorJSON(w, r, http.StatusBadRequest, err, "invalid credential fetch request body", 0)
			return
		}

		_, err = core.IDFromString(basicMessage.To)
		if err != nil {
			rest.ErrorJSON(w, r, http.StatusBadRequest, err, "invalid issuer id in base message", 0)
			return
		}

		var claimID uuid.UUID
		claimID, err = uuid.Parse(fetchRequestBody.ID)
		if err != nil {
			rest.ErrorJSON(w, r, http.StatusBadRequest, err, "invalid claim id in fetch request body", 0)
			return
		}

		var claim *models.Claim
		claimIDByte := []byte(claimID.String())
		claim, err = comm.claimService.GetByID(claimIDByte)

		if err != nil {
			rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "failed get claim by claimID", 0)
			return
		}

		if claim.OtherIdentifier != basicMessage.From {
			rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "claim doesn't relate to sender", 0)
			return
		}

		var cred *verifiable.Iden3Credential
		cred, err = comm.schemaService.FromClaimToIden3Credential(*claim)
		if err != nil {
			rest.ErrorJSON(w, r, http.StatusInternalServerError, err, "failed convert claim", 0)
			return
		}

		resp := protocol.CredentialIssuanceMessage{
			ID:       uuid.NewString(),
			Type:     protocol.CredentialIssuanceResponseMessageType,
			ThreadID: basicMessage.ThreadID,
			Body:     protocol.IssuanceMessageBody{Credential: *cred},
			From:     basicMessage.To,
			To:       basicMessage.From,
		}

		if err != nil {
			rest.ErrorJSON(w, r, http.StatusBadRequest, err, "failed marshal issuance response", 0)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, resp)
	}

	return
}

// Unpack returns unpacked message from transport envelope with verification of zeroknowledge proof
func Unpack(circuitPath string, envelope []byte) (*iden3comm.BasicMessage, error) {

	token, err := jwz.Parse(string(envelope))
	if err != nil {
		return nil, err
	}

	verificationKey, err := getVerificationKey(circuitPath, circuits.AuthCircuitID)
	if err != nil {
		return nil, errors.New("message was packed with unsupported circuit")
	}

	isValid, err := token.Verify(verificationKey)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errors.New("message proof is invalid")
	}

	var msg iden3comm.BasicMessage
	err = json.Unmarshal(token.GetPayload(), &msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &msg, err
}

func getVerificationKey(basePath string, circuitID circuits.CircuitID) ([]byte, error) {
	return getPathToFile(basePath, circuitID, "verification_key.json")
}

func getPathToFile(basePath string, circuitID circuits.CircuitID, fileName string) ([]byte, error) {
	path := filepath.Join(basePath, string(circuitID), fileName)
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed open file '%s' by path '%s'", fileName, path)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed read file '%s' by path '%s'", fileName, path)
	}
	return data, nil
}
