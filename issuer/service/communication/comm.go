package communication

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/iden3/go-circuits"
	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
	"github.com/iden3/iden3comm/protocol"
	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"issuer/service/cfgs"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// TODO:
// 1. Update the new key dir direction to the new location
// 2. Update the url to callback, only issuer url is relevant

var requestTracker = cache.New(60*time.Minute, 60*time.Minute)

func NewCommunicationHandler(issuerId string, cfg cfgs.IssuerConfig) *Handler {
	return &Handler{
		issuerId:  issuerId,
		keyDir:    cfg.KeyDir,
		publicUrl: cfg.PublicUrl,
	}
}

type Handler struct {
	keyDir    string
	publicUrl string
	issuerId  string
}

// sending sign in request to the client (move it to the issuer communication (identity))
func (h *Handler) GetAuthRequest(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler.GetAuthRequest() invoked")

	sessionID := rand.Intn(1000000)

	hostUrl := h.publicUrl
	if len(hostUrl) == 0 {
		log.Fatal("host-url is not set")
	}
	t := r.URL.Query().Get("type")
	uri := fmt.Sprintf("%s/api/v1/callback?sessionId=%s", hostUrl, strconv.Itoa(sessionID))

	var request protocol.AuthorizationRequestMessage
	request = auth.CreateAuthorizationRequestWithMessage("test flow", "message to sign", "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ", uri)

	request.ID = "7f38a193-0918-4a48-9fac-36adfdb8b542"
	request.ThreadID = "7f38a193-0918-4a48-9fac-36adfdb8b542"

	if t != "" {
		var mtpProofRequest protocol.ZeroKnowledgeProofRequest
		mtpProofRequest.ID = 1
		mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigCircuitID)
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
	}

	sId := strconv.Itoa(sessionID)
	requestTracker.Set(sId, request, cache.DefaultExpiration)

	msgBytes, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "x-id")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-id", strconv.Itoa(sessionID))
	w.WriteHeader(http.StatusOK)

	w.Write(msgBytes)
	return
}

func (h *Handler) GetAgeClaimOffer(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler.GetAgeClaimOffer() invoked")

	userId := chi.URLParam(r, "user-id")
	claimId := chi.URLParam(r, "claim-id")

	res := &protocol.CredentialsOfferMessage{
		ID:       "7f38a193-0918-4a48-9fac-36adfdb8b542",
		Typ:      "application/iden3comm-plain-json",
		Type:     "https://iden3-communication.io/credentials/1.0/offer",
		ThreadID: "f7a3fae9-ecf1-4603-804e-8ff1c7632636",
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(msgBytes)
	return
}

// currently in the frontend but need to move here
//func (h *Handler) getOfferClaimRequest() (*protocol.CredentialsOfferMessage, error) {
//	logger.Debug("Handler.getOfferClaimRequest() invoked")
//
//	return &protocol.CredentialsOfferMessage{
//		ID:       "7f38a193-0918-4a48-9fac-36adfdb8b542",
//		Typ:      "application/iden3comm-plain-json",
//		Type:     "https://iden3-communication.io/credentials/1.0/offer",
//		ThreadID: "f7a3fae9-ecf1-4603-804e-8ff1c7632636",
//		Body: protocol.CredentialsOfferMessageBody{
//			URL: h.cfg.PublicUrl + `/api/v1/agent`,
//			Credentials: []protocol.CredentialOffer{
//				protocol.CredentialOffer{
//					ID:          "7f38a193-0918-4a48-", // mock value - need to fix this (claim id in the front end)
//					Description: "KYCAgeCredential",
//				},
//			}},
//		From: "123214512", // mock value - need to fix this
//		To:   "issuerId",  // mock value - need to fix this
//	}, nil
//}

// Handle the sign in response from the user
// Callback works with sign-in callbacks
func (h *Handler) Callback(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler.Callback() invoked")

	switch r.Method {
	case http.MethodPost:
		// get query params from request
		sessionID := r.URL.Query().Get("sessionId")
		tokenBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cant read bytes", http.StatusBadRequest)
			return
		}

		keyDIR := h.keyDir
		if len(keyDIR) == 0 {
			log.Fatal("host-url is not set")
		}

		resolver := state.ETHResolver{
			RPCUrl:   "https://polygon-mumbai.g.alchemy.com/v2/qe4__xaBZRG7AN0AuKrw72C7REjSt1DM",
			Contract: "0x46Fd04eEa588a3EA7e9F055dd691C688c4148ab3",
		}

		var verificationKeyloader = &loaders.FSKeyLoader{Dir: keyDIR}
		verifier := auth.NewVerifier(verificationKeyloader, loaders.DefaultSchemaLoader{IpfsURL: "ipfs.io"}, resolver)
		authRequest, wasFound := requestTracker.Get(sessionID)
		if wasFound == false {
			logger.Errorf("auth request was not found for session ID:", sessionID)
			http.Error(w, "no request", http.StatusBadRequest)
			return
		}

		arm, err := verifier.FullVerify(r.Context(), string(tokenBytes), authRequest.(protocol.AuthorizationRequestMessage))
		if err != nil { // if enter here, the verification is result is false
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		m := make(map[string]interface{})

		m["id"] = arm.From

		mBytes, err := json.Marshal(m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO: Where is this being used in production?
		requestTracker.Set(sessionID, m, cache.DefaultExpiration)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mBytes)

		return
	default:
		return
	}
}

// GetRequestStatus checks response status
func (h *Handler) GetRequestStatus(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Handler.GetRequestStatus() invoked")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	logger.Tracef("cache - getting id %s from cahce\n", id)
	item, ok := requestTracker.Get(id)
	if !ok {
		logger.Errorf("item not found %v", id)
		http.NotFound(w, r)
		return
	}

	switch item.(type) {
	case protocol.AuthorizationRequestMessage:
		http.Error(w, errors.New("no authorization response yet").Error(), http.StatusBadRequest)
		return
	case map[string]interface{}:
		b, err := json.Marshal(item)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}
}
