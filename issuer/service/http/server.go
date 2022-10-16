package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	logger "github.com/sirupsen/logrus"
	"issuer/service/claims"
	"issuer/service/identity"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	address    string
	identity   *identity.Identity
	claims     *claims.Claims
}

func NewServer(serviceAddress string, identity *identity.Identity, claims *claims.Claims) *Server {
	return &Server{
		address:  serviceAddress,
		identity: identity,
		claims:   claims,
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
			//rclaim.Get("/revocation/status/{nonce}", s.Claim.RevocationStatus)
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
	res := &GetIdentityResponse{Identity: s.identity.GetIdentity()}

	EncodeResponse(w, 200, res)
}

func (s *Server) getClaim(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) createClaim(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) getRevocationStatus(w http.ResponseWriter, r *http.Request) {

	return
}

func (s *Server) getAgent(w http.ResponseWriter, r *http.Request) {

	return
}


