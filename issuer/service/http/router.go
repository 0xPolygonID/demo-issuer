package http

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	logger "github.com/sirupsen/logrus"
)

func newRouter(s *Server) chi.Router {
	logger.Debug("Server.newRouter() invoked")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	corsMiddleware := cors.New(cors.Options{
		//AllowedOrigins:   []string{"localhost", "127.0.0.1", "*"},
		AllowedOrigins:   []string{"http://*", "https://*", "*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)

	r.Route("/api/v1", func(root chi.Router) {
		root.Use(render.SetContentType(render.ContentTypeJSON))

		root.Route("/identity", func(r chi.Router) {
			r.Post("/", s.getIdentity)
			r.Post("/publish", s.publish)
		})

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

		root.Route("/age-verification-request", func(agent chi.Router) {
			agent.Get("/", s.issuer.CommHandler.GetAgeVerificationRequest)
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
