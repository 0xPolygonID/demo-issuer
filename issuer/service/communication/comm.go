package communication

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/iden3/go-circuits"
	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
	"github.com/iden3/iden3comm/protocol"
	"github.com/patrickmn/go-cache"
	logger "github.com/sirupsen/logrus"
	"issuer/service/cfgs"
	"math/rand"
	"strconv"
	"time"
)

var userSessionTracker = cache.New(60*time.Minute, 60*time.Minute)

func NewCommunicationHandler(issuerId string, cfg *cfgs.IssuerConfig) *Handler {
	return &Handler{
		issuerId:   issuerId,
		keyDir:     cfg.CircuitsDir,
		publicUrl:  cfg.PublicUrl,
		NodeRpcUrl: cfg.NodeRpcUrl,
		ipfsUrl:    cfg.IpfsUrl,
	}
}

type Handler struct {
	keyDir     string
	publicUrl  string
	issuerId   string
	NodeRpcUrl string
	ipfsUrl    string
}

func (h *Handler) GetAuthVerificationRequest() ([]byte, string, error) {
	logger.Debug("Communication.GetAuthVerificationRequest() invoked")

	sId := strconv.Itoa(rand.Intn(1000000))
	uri := fmt.Sprintf("%s/api/v1/callback?sessionId=%s", h.publicUrl, sId)

	var request protocol.AuthorizationRequestMessage
	request = auth.CreateAuthorizationRequestWithMessage("test flow", "message to sign", "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ", uri)

	request.ID = uuid.New().String()
	request.ThreadID = uuid.New().String()

	userSessionTracker.Set(sId, request, cache.DefaultExpiration)

	msgBytes, err := json.Marshal(request)
	if err != nil {
		return nil, "", fmt.Errorf("error marshalizing response: %v", err)
	}

	return msgBytes, sId, nil
}

func verifySupportedAgeCircuit(circuitType string) error {
	if circuitType == string(circuits.AtomicQueryMTPCircuitID) ||
		circuitType == string(circuits.AtomicQuerySigCircuitID) {
		return nil
	}

	return fmt.Errorf("'%s' unsupported circuit type")
}

func (h *Handler) GetAgeVerificationRequest(circuitType string) ([]byte, string, error) {
	logger.Debug("Communication.GetAgeVerificationRequest() invoked")

	sId := strconv.Itoa(rand.Intn(1000000))

	hostUrl := h.publicUrl
	if len(hostUrl) == 0 {
		return nil, "", fmt.Errorf("host-url is not set")
	}
	uri := fmt.Sprintf("%s/api/v1/callback?sessionId=%s", hostUrl, sId)

	var request protocol.AuthorizationRequestMessage
	request = auth.CreateAuthorizationRequestWithMessage("test flow", "message to sign", "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ", uri)

	request.ID = uuid.New().String()
	request.ThreadID = uuid.New().String()

	if err := verifySupportedAgeCircuit(circuitType); err != nil {
		return nil, "", err
	}

	var mtpProofRequest protocol.ZeroKnowledgeProofRequest
	mtpProofRequest.ID = 1
	mtpProofRequest.CircuitID = circuitType
	mtpProofRequest.Rules = map[string]interface{}{
		"query": pubsignals.Query{
			AllowedIssuers: []string{"*"},
			Req: map[string]interface{}{
				"birthday": map[string]interface{}{
					"$lt": 20000101,
				},
			},
			Schema: protocol.Schema{
				URL:  "https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/kyc-v2.json-ld",
				Type: "KYCAgeCredential",
			},
		},
	}
	request.Body.Scope = append(request.Body.Scope, mtpProofRequest)

	userSessionTracker.Set(sId, request, cache.DefaultExpiration)

	msgBytes, err := json.Marshal(request)
	if err != nil {
		return nil, "", fmt.Errorf("error marshalizing response: %v", err)
	}

	return msgBytes, sId, nil
}

func (h *Handler) GetAgeClaimOffer(userId, claimId string) ([]byte, error) {
	logger.Debug("Communication.GetAgeClaimOffer() invoked")

	res := &protocol.CredentialsOfferMessage{
		ID:       uuid.New().String(),
		Typ:      "application/iden3comm-plain-json",
		Type:     "https://iden3-communication.io/credentials/1.0/offer",
		ThreadID: uuid.New().String(),
		Body: protocol.CredentialsOfferMessageBody{
			URL: h.publicUrl + `/api/v1/agent`,
			Credentials: []protocol.CredentialOffer{
				protocol.CredentialOffer{
					ID:          claimId,
					Description: "KYCAgeCredential",
				},
			}},
		From: userId,
		To:   h.issuerId,
	}

	msgBytes, err := json.Marshal(res)
	if err != nil {
		return nil, fmt.Errorf("error marshalizing response: %v", err)
	}

	return msgBytes, nil
}

func (h *Handler) Callback(sId string, tokenBytes []byte) ([]byte, error) {
	logger.Debug("Communication.Callback() invoked")

	authRequest, wasFound := userSessionTracker.Get(sId)
	if wasFound == false {
		err := fmt.Errorf("auth request was not found for session ID: %s", sId)
		logger.Errorf(err.Error())
		return nil, err
	}

	resolver := state.ETHResolver{
		RPCUrl:   h.NodeRpcUrl,
		Contract: "0x46Fd04eEa588a3EA7e9F055dd691C688c4148ab3",
	}

	var verificationKeyLoader = &loaders.FSKeyLoader{Dir: h.keyDir}
	verifier := auth.NewVerifier(verificationKeyLoader, loaders.DefaultSchemaLoader{IpfsURL: h.ipfsUrl}, resolver)

	arm, err := verifier.FullVerify(context.Background(), string(tokenBytes), authRequest.(protocol.AuthorizationRequestMessage))
	if err != nil { // the verification result is false
		return nil, err
	}

	m := make(map[string]interface{})
	m["id"] = arm.From

	mBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("error marshalizing response: %v", err)
	}

	userSessionTracker.Set(sId, m, cache.DefaultExpiration)

	return mBytes, nil
}

// GetRequestStatus checks response status
func (h *Handler) GetRequestStatus(id string) ([]byte, error) {
	logger.Debug("Communication.Callback() invoked")

	logger.Tracef("cache - getting id %s from cahce\n", id)
	item, ok := userSessionTracker.Get(id)
	if !ok {
		logger.Tracef("item not found %v", id)
		return nil, nil
	}

	switch item.(type) {
	case protocol.AuthorizationRequestMessage:
		logger.Warn("no authorization response yet - no data available for this request")
		return nil, nil
	case map[string]interface{}:
		b, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("error marshalizing response: %v", err)
		}
		return b, nil
	}

	return nil, fmt.Errorf("unknown item return from tracker (type %T)", item)
}
