package command

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iden3/go-iden3-core"
	"github.com/iden3/go-jwz"
	"github.com/iden3/iden3comm"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"issuer/service/claim"
	"issuer/service/identity/state"
	"issuer/utils"
)

type Handler struct {
	issuerIdentifier *core.ID
	idenState        *state.IdentityState
	keysPath         string
}

func NewHandler(
	issuerIdentifier *core.ID,
	state *state.IdentityState,
	keysPath string,
) *Handler {
	return &Handler{
		issuerIdentifier: issuerIdentifier,
		idenState:        state,
		keysPath:         keysPath,
	}
}

// Handle POST /api/v1/agent
func (comm *Handler) Handle(body []byte) (*protocol.CredentialIssuanceMessage, error) {
	logger.Debug("CommandHandler.Handle() invoked")

	basicMessage, err := comm.unpack(body)
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

	claimID, err := uuid.Parse(fetchRequestBody.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid claim id in fetch request body, err: %v", err)
	}

	c, err := comm.idenState.Claims.GetClaim([]byte(claimID.String()))
	if err != nil {
		return nil, err
	}
	if c.OtherIdentifier != basicMessage.From {
		return nil, err
	}

	if !comm.idenState.CommittedState.IsLatestStateGenesis {
		claimIdx, err := c.CoreClaim.HIndex()
		if err != nil {
			return nil, err
		}
		mtp, err := comm.idenState.GetMTPProof(comm.issuerIdentifier, claimIdx)
		if err != nil {
			return nil, err
		}
		c.MTPProof, err = json.Marshal(mtp)
		if err != nil {
			return nil, err
		}
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

// unpack returns unpacked message from transport envelope with verification of zeroknowledge proof
func (comm *Handler) unpack(envelope []byte) (*iden3comm.BasicMessage, error) {
	logger.Debug("CommandHandler.unpack() invoked")

	token, err := jwz.Parse(string(envelope))
	if err != nil {
		return nil, err
	}

	verificationKey, err := utils.ReadFileByPath(comm.keysPath, "auth.json")
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
