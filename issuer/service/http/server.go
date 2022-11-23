package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	logger "github.com/sirupsen/logrus"
	"io"
	"issuer/service/identity"
	"issuer/service/models"
	"net/http"
	"strconv"
)

type Server struct {
	httpServer *http.Server
	address    string
	issuer     *identity.Identity
}

func NewServer(localHostAdd string, issuer *identity.Identity) *Server {

	return &Server{
		address: localHostAdd,
		issuer:  issuer,
	}
}

func (s *Server) Close(ctx context.Context) error {
	logger.Debug("Server.Close() invoked")

	if s.httpServer == nil {
		return fmt.Errorf("error on server.Close(). http server is not running")
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Run() error {
	logger.Debug("Server.Run() invoked")

	logger.Infof("starting HTTP server (address: %s)", s.address)

	s.httpServer = &http.Server{Addr: s.address, Handler: newRouter(s)}
	return s.httpServer.ListenAndServe()
}

func (s *Server) getIdentity(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getIdentity() invoked")

	iden, err := s.issuer.GetIdentity()
	if err != nil {
		logger.Errorf("Server -> issuer.GetIdentity() return err, err: %v", err)
		EncodeResponse(w, 500, err)
	}

	EncodeResponse(w, 200, iden)
}

func (s *Server) createClaim(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.createClaim() invoked")

	req := &models.CreateClaimRequest{}
	if err := JsonToStruct(r, req); err != nil {
		logger.Error("cannot unmarshal json body, err: %v", err)
		EncodeResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := s.issuer.CreateClaim(req)
	if err != nil {
		logger.Errorf("Server -> issuer.CreateClaim() return err, err: %v", err)
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("can't parse claim id param - %v", err))
		return
	}

	EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) getClaim(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getClaim() invoked")

	claimID := chi.URLParam(r, "id")

	if claimID == "" {
		logger.Errorf("Server error on getClaim(), claimID is emppty")
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("no claim identifier - can't  parse claim id param"))
		return
	}

	res, err := s.issuer.GetClaim(claimID)
	if err != nil {
		logger.Errorf("Server -> issuer.CreateClaim() return err, err: %v", err)
		EncodeResponse(w, http.StatusNotFound, fmt.Errorf("can't get claim %s, err: %v", claimID, err))
		return
	}

	EncodeResponse(w, 200, res)
}

func (s *Server) getRevocationStatus(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getRevocationStatus() invoked")

	nonce, err := strconv.ParseUint(chi.URLParam(r, "nonce"), 10, 64)
	if err != nil {
		logger.Errorf("error on parsing nonce, err: %v", err)
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("error on parsing nonce input"))
		return
	}

	res, err := s.issuer.GetRevocationStatus(nonce)
	if err != nil {
		logger.Errorf("Server -> issuer.GetRevocationStatus() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't generate non revocation proof for revocation nonce: %d. err: %v", nonce, err))
		return
	}
	EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.callback() invoked")

	sessionID := r.URL.Query().Get("sessionId")
	tokenBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Server.callback() error reading request body, err: %v", err)
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("can't read request body"))
		return
	}

	resB, err := s.issuer.CommHandler.Callback(sessionID, tokenBytes)
	if err != nil {
		logger.Errorf("Server.callback() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Errorf("can't handle callback request"))
		return
	}

	EncodeByteResponse(w, http.StatusOK, resB)
}

func (s *Server) getRequestStatus(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getRequestStatus() invoked")

	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Errorf("Server.getRequestStatus() url parameter has invalid values")
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("url parameter has invalid values"))
		return
	}

	resB, err := s.issuer.CommHandler.GetRequestStatus(id)
	if err != nil {
		logger.Errorf("Server -> issuer.CommHandler.GetRequestStatus() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't get request status. err: %v", err))
		return
	}

	if resB == nil {
		EncodeResponse(w, http.StatusNotFound, fmt.Errorf("can't get request status with id: %s", id))
		return
	}

	EncodeByteResponse(w, http.StatusOK, resB)
}

func (s *Server) getAgeClaimOffer(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getAgeClaimOffer() invoked")

	userId := chi.URLParam(r, "user-id")
	claimId := chi.URLParam(r, "claim-id")
	if userId == "" || claimId == "" {
		logger.Errorf("Server.getAgeClaimOffer() url parameters has invalid values")
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("url parameters has invalid values"))
		return
	}

	resB, err := s.issuer.CommHandler.GetAgeClaimOffer(userId, claimId)
	if err != nil {
		logger.Errorf("Server -> issuer.CommHandler.getAgeClaimOffer() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't get age claim offer. err: %v", err))
		return
	}

	EncodeByteResponse(w, http.StatusOK, resB)
}

func (s *Server) getAuthVerificationRequest(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getAuthVerificationRequest() invoked")

	resB, sessionId, err := s.issuer.CommHandler.GetAuthVerificationRequest()
	if err != nil {
		logger.Errorf("Server -> issuer.CommHandler.GetAuthVerificationRequest() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't get auth verification request. err: %v", err))
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "x-id")
	w.Header().Set("x-id", sessionId)
	EncodeByteResponse(w, http.StatusOK, resB)
}

func (s *Server) getAgeVerificationRequest(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getAgeVerificationRequest() invoked")

	circuitType := r.URL.Query().Get("circuitType")
	if circuitType == "" {
		logger.Errorf("Circuit type is empty")
		EncodeResponse(w, http.StatusInternalServerError, "Circuit type is empty")
		return
	}

	resB, sessionId, err := s.issuer.CommHandler.GetAgeVerificationRequest(circuitType)
	if err != nil {
		logger.Errorf("Server -> issuer.CommHandler.GetAuthVerificationRequest() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't get auth verification request. err: %v", err))
		return
	}

	w.Header().Set("Access-Control-Expose-Headers", "x-id")
	w.Header().Set("x-id", sessionId)
	EncodeByteResponse(w, http.StatusOK, resB)
}

func (s *Server) agent(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.agent() invoked")

	bodyB, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Server error on parsing request body (agent), err: %v", err)
		EncodeResponse(w, http.StatusBadRequest, "can't parse request, err: "+err.Error())
		return
	}

	res, err := s.issuer.CmdHandler.Handle(bodyB)
	if err != nil {
		logger.Errorf("Server -> commHandler.Handle() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, "error on handle income request, err: "+err.Error())
		return
	}

	EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) publish(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.publish() invoked")

	txHex, err := s.issuer.PublishLatestState(r.Context())
	if err != nil {
		logger.Errorf("Server -> issuer.publish() return err, err: %v", err)
		EncodeResponse(w, http.StatusInternalServerError, "error on publishing latest state: "+err.Error())
		return
	}

	res := struct {
		Hex string `json:"hex"`
	}{Hex: txHex}
	EncodeResponse(w, http.StatusOK, res)
}
