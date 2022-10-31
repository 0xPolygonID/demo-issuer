package command

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-core"
	"github.com/iden3/go-jwz"
	"github.com/iden3/iden3comm"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"issuer/service/claim"
	"issuer/service/identity/state"
	"os"
	"path/filepath"
)

// MaxBodySizeBytes defines is 2 MB for protocol message
const MaxBodySizeBytes = 2 * 1000 * 1000

// Handler provides handler for zk proof
type Handler struct {
	idenState   *state.IdentityState
	circuitPath string
}

// NewIDEN3CommHandler inits IDEN3Comm handler
func NewHandler(
	state *state.IdentityState,
	circuitPath string,
) *Handler {
	return &Handler{
		circuitPath: circuitPath,
		idenState:   state,
	}
}

// Handle POST /api/v1/agent
func (comm *Handler) Handle(body []byte) (*protocol.CredentialIssuanceMessage, error) {
	logger.Debug("CommandHandler.Handle() invoked")

	basicMessage, err := Unpack(comm.circuitPath, body)
	if err != nil {
		return nil, err
	}

	if basicMessage.Type != protocol.CredentialFetchRequestMessageType {
		return nil, fmt.Errorf("unsupported protocol message type")
	}

	fetchRequestBody := &protocol.CredentialFetchRequestMessageBody{}
	err = json.Unmarshal(basicMessage.Body, fetchRequestBody)
	if err != nil {
		return nil, err
	}

	_, err = core.IDFromString(basicMessage.To)
	if err != nil {
		return nil, err
	}

	c, err := comm.idenState.Claims.GetClaim([]byte(fetchRequestBody.ID))
	if err != nil {
		return nil, err
	}
	if c.OtherIdentifier != basicMessage.From {
		return nil, err
	}

	cred, err := claim.ClaimModelToIden3Credential(c)
	if err != nil {
		return nil, err
	}

	resp := &protocol.CredentialIssuanceMessage{
		ID:       uuid.NewString(),
		Type:     protocol.CredentialIssuanceResponseMessageType,
		ThreadID: basicMessage.ThreadID,
		Body:     protocol.IssuanceMessageBody{Credential: *cred},
		From:     basicMessage.To,
		To:       basicMessage.From,
	}

	return resp, nil
}

// Unpack returns unpacked message from transport envelope with verification of zeroknowledge proof
func Unpack(circuitPath string, envelope []byte) (*iden3comm.BasicMessage, error) {
	logger.Debug("CommandHandler.Unpack() invoked")

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

	msg := &iden3comm.BasicMessage{}
	err = json.Unmarshal(token.GetPayload(), msg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if msg.From == "" {
		return nil, fmt.Errorf("invalid message: from field is invalid")
	}

	if msg.To == "" {
		return nil, fmt.Errorf("invalid message: to field is invalid")
	}

	return msg, err
}

func getVerificationKey(basePath string, circuitID circuits.CircuitID) ([]byte, error) {
	logger.Debug("CommandHandler.getVerificationKey() invoked")

	return getPathToFile(basePath, circuitID, "verification_key.json")
}

func getPathToFile(basePath string, circuitID circuits.CircuitID, fileName string) ([]byte, error) {
	logger.Debug("CommandHandler.getPathToFile() invoked")

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
