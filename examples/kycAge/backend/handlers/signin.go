package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	//auth "github.com/bajpai244/go-iden3-auth"
	"github.com/iden3/go-circuits"
	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
	"github.com/iden3/iden3comm/protocol"
	"github.com/patrickmn/go-cache"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var cacheStorage = cache.New(60*time.Minute, 60*time.Minute)

func GetQR(w http.ResponseWriter, r *http.Request) {
	sessionID := rand.Intn(1000000)

	hostUrl := os.Getenv("HOST_URL")
	t := r.URL.Query().Get("type")
	uri := fmt.Sprintf("%s/api/callback?sessionId=%s", hostUrl, strconv.Itoa(sessionID))

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

	cacheStorage.Set(strconv.Itoa(sessionID), request, cache.DefaultExpiration)

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

// Callback works with sign-in callbacks
func Callback(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// get query params from request
		sessionID := r.URL.Query().Get("sessionId")
		keyDIR := os.Getenv("KEY_DIR")
		tokenBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cant read bytes", http.StatusBadRequest)
			return
		}

		if keyDIR == "" {
			fmt.Println("please, provide KEY_DIR env")
			http.Error(w, "Verifier can't verify token. Key loader is not setup", http.StatusInternalServerError)
			return
		}

		resolver := state.ETHResolver{
			RPCUrl:   "https://polygon-mumbai.g.alchemy.com/v2/qe4__xaBZRG7AN0AuKrw72C7REjSt1DM",
			Contract: "0x46Fd04eEa588a3EA7e9F055dd691C688c4148ab3",
		}

		var verificationKeyloader = &loaders.FSKeyLoader{Dir: keyDIR}
		verifier := auth.NewVerifier(verificationKeyloader, loaders.DefaultSchemaLoader{IpfsURL: "ipfs.io"}, resolver)
		authRequest, wasFound := cacheStorage.Get(sessionID)
		if wasFound == false {
			fmt.Println("auth request was not found for session ID:", sessionID)
			http.Error(w, "no request", http.StatusBadRequest)
			return
		}

		arm, err := verifier.FullVerify(r.Context(), string(tokenBytes),
			authRequest.(protocol.AuthorizationRequestMessage))
		if err != nil {
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
		cacheStorage.Set(sessionID, m, cache.DefaultExpiration)

		//nolint // -
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(mBytes)

		return
	default:
		return
	}
}

// Status checks response status
func Status(w http.ResponseWriter, r *http.Request) {
	//nolint // -
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")

	item, ok := cacheStorage.Get(id)

	if !ok {
		log.Printf("item not found %v", id)
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