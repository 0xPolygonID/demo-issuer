package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
	"issuer/service/contract"
	"issuer/service/identity"
	"net/http"
	"strconv"
)

type Server struct {
	httpServer *http.Server
	address    string
	issuer     *identity.Identity
}

func NewServer(serviceAddress string, issuer *identity.Identity) *Server {
	return &Server{
		address: serviceAddress,
		issuer:  issuer,
	}
}

func (s *Server) Close(ctx context.Context) error {
	if s.httpServer == nil {
		return fmt.Errorf("error on server.Close(). http server is not running")
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Run() error {
	logger.Info("starting HTTP server (address: %s)", s.address)

	s.httpServer = &http.Server{Addr: s.address, Handler: s.newRouter()}
	return s.httpServer.ListenAndServe()
}

func (s *Server) newRouter() chi.Router {
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
			//rclaim.Get("/{id}", s.Claim.GetClaimByID)
			claims.Get("/{id}", s.getClaim)

			//rclaim.Post("/", s.Claim.CreateClaim)
			claims.Post("/", s.createClaim)
		})

		root.Route("/revocations/{nonce}", func(revs chi.Router) {
			//rclaim.Get("/revocation/status/{nonce}", s.Claim.GetRevocationStatusResponse)
			revs.Post("/", s.getRevocationStatus)
		})

		root.Route("/agent", func(agent chi.Router) {
			//internal.Post("/", s.IDEN3Comm.Handle)
			agent.Post("/", s.getAgent)
		})
	})

	return r
}

func (s *Server) getIdentity(w http.ResponseWriter, r *http.Request) {
	iden, err := s.issuer.GetIdentity()
	if err != nil {
		EncodeResponse(w, 500, err)
	}

	//res := &contract.GetIdentityResponse{
	//	Identifier: iden.Identifier,
	//	State: &contract.IdentityState{
	//		Identifier:         iden.Identifier,
	//		StateID:            iden.StateID,
	//		State:              iden.State,
	//		ClaimsTreeRoot:     iden.ClaimsTreeRoot,
	//		RevocationTreeRoot: iden.RevocationTreeRoot,
	//		RootOfRoots:        iden.RootOfRoots,
	//	},
	//}

	EncodeResponse(w, 200, iden)
}

func (s *Server) createClaim(w http.ResponseWriter, r *http.Request) {
	req := &contract.CreateClaimRequest{}
	if err := JsonToStruct(r, req); err != nil {
		logger.Error("cannot unmarshal json body")
		EncodeResponse(w, http.StatusBadRequest, err)
		return
	}

	// validate

	// convert to biz-model

	// call issuer add claim
	res, err := s.issuer.CreateClaim(req)
	if err != nil {
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("can't parse claim id param - %v", err))
		return
	}

	EncodeResponse(w, http.StatusOK, res)
}

func (s *Server) getClaim(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	if idParam == "" {
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("no claim identifier - can't  parse claim id param"))
		return
	}

	claimID, err := uuid.Parse(idParam)
	if err != nil {
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("can't parse claim id param - %v", err))
		return
	}

	claim, err := s.issuer.GetClaim(claimID.String())
	if err != nil {
		EncodeResponse(w, http.StatusNotFound, fmt.Errorf("can't get claim %s, err: %v", claimID.String(), err))
		return
	}

	res := &contract.GetClaimResponse{
		ID:   string(claim.ID),
		Type: claim.SchemaType,
	}

	EncodeResponse(w, 200, res)
}

func (s *Server) getRevocationStatus(w http.ResponseWriter, r *http.Request) {

	nonce, err := strconv.ParseUint(chi.URLParam(r, "nonce"), 10, 64)
	if err != nil {
		EncodeResponse(w, http.StatusBadRequest, fmt.Errorf("error on parsing nonce input"))
	}

	proof, err := s.issuer.GetRevocation(nonce)
	if err != nil {
		EncodeResponse(w, http.StatusInternalServerError, fmt.Sprintf("can't generate non revocation proof for revocation nonce: %d. err: %v", nonce, err))
		return
	}
	EncodeResponse(w, http.StatusOK, proof)
}

func (s *Server) getAgent(w http.ResponseWriter, r *http.Request) {

	return
}
