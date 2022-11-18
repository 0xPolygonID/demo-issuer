package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	logger "github.com/sirupsen/logrus"
	"io"
	http2 "issuer/http"
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

	s.httpServer = &http.Server{Addr: s.address, Handler: s.newRouter()}
	return s.httpServer.ListenAndServe()
}

func (s *Server) newRouter() chi.Router {
	logger.Debug("Server.newRouter() invoked")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost", "127.0.0.1", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)

	r.Route("/api/v1", func(root chi.Router) {
		root.Use(render.SetContentType(render.ContentTypeJSON))

		root.Get("/identity", s.getIdentity)

		root.Route("/claims", func(claims chi.Router) {
			claims.Get("/{id}", s.getClaim)
			claims.Post("/", s.createClaim)

			claims.Route("/offers", func(claimRequests chi.Router) {
				claimRequests.Get("/{user-id}/{claim-id}", s.issuer.CommHandler.GetAgeClaimOffer)
			})

			claims.Route("/revocations", func(revs chi.Router) {
				revs.Get("/{nonce}", s.getRevocationStatus)
			})

		})

		root.Route("/agent", func(agent chi.Router) {
			agent.Post("/", s.getAgent)
		})

		root.Route("/sign-in", func(agent chi.Router) {
			agent.Get("/", s.issuer.CommHandler.GetAuthVerificationRequest)
		})
		root.Route("/callback", func(agent chi.Router) {
			agent.Post("/", s.issuer.CommHandler.Callback)
		})
		root.Route("/status", func(agent chi.Router) {
			agent.Get("/", s.issuer.CommHandler.GetRequestStatus)
		})

	})

	return r
}

func (s *Server) getIdentity(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getIdentity() invoked")

	iden, err := s.issuer.GetIdentity()
	if err != nil {
		logger.Errorf("Server -> issuer.GetIdentity() return err, err: %v", err)
		http2.EncodeResponse(w, 500, err)
	}

	http2.EncodeResponse(w, 200, iden)
}

func (s *Server) createClaim(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.createClaim() invoked")

	req := &models.CreateClaimRequest{}
	if err := http2.JsonToStruct(r, req); err != nil {
		logger.Error("cannot unmarshal json body, err: %v", err)
		http2.EncodeResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate

	// convert to biz-model

	// call issuer add claim
	res, err := s.issuer.CreateClaim(req)
	if err != nil {
		logger.Errorf("Server -> issuer.CreateClaim() return err, err: %v", err)
		http2.EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("can't parse claim id param - %v", err))
		return
	}

	http2.EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) getClaim(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getClaim() invoked")

	claimID := chi.URLParam(r, "id")

	if claimID == "" {
		logger.Errorf("Server error on getClaim(), claimID is emppty")
		http2.EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("no claim identifier - can't  parse claim id param"))
		return
	}

	res, err := s.issuer.GetClaim(claimID)
	if err != nil {
		logger.Errorf("Server -> issuer.CreateClaim() return err, err: %v", err)
		http2.EncodeResponse(w, http.StatusNotFound, fmt.Errorf("can't get claim %s, err: %v", claimID, err))
		return
	}

	http2.EncodeResponse(w, 200, res)
}

func (s *Server) getRevocationStatus(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getRevocationStatus() invoked")

	nonce, err := strconv.ParseUint(chi.URLParam(r, "nonce"), 10, 64)
	if err != nil {
		logger.Errorf("Server error on parsing nonce, err: %v", err)
		http2.EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("error on parsing nonce input"))
		return
	}

	res, err := s.issuer.GetRevocationStatus(nonce)
	if err != nil {
		logger.Errorf("Server -> issuer.GetRevocationStatus() return err, err: %v", err)
		http2.EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't generate non revocation proof for revocation nonce: %d. err: %v", nonce, err))
		return
	}
	http2.EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) getAgent(w http.ResponseWriter, r *http.Request) {
	logger.Debug("Server.getAgent() invoked")

	bodyB, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("Server error on parsing request body (getAgent), err: %v", err)
		http2.EncodeResponse(w, http.StatusBadRequest, "can't parse request, err: "+err.Error())
		return
	}

	res, err := s.issuer.CmdHandler.Handle(bodyB)
	if err != nil {
		logger.Errorf("Server -> commHandler.Handle() return err, err: %v", err)
		http2.EncodeResponse(w, http.StatusInternalServerError, "error on handle income request, err: "+err.Error())
		return
	}

	http2.EncodeResponse(w, http.StatusOK, res)
}
