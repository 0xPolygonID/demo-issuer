package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"lightissuer/app/handlers"
	"net/http"
)

// Handlers server handlers
type Handlers struct {
	/* Put handlers here*/
	Identity  *handlers.IdentityHandler
	Claim     *handlers.ClaimHandler
	IDEN3Comm *handlers.IDEN3CommHandler
}

// Routes chi http routs configuration
func (s *Handlers) Routes() chi.Router {
	r := chi.NewRouter()

	// Basic CORS
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"localhost", "127.0.0.1", "*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type",
			"X-CSRF-Token"},
		AllowCredentials: true,
	})

	r.Use(corsMiddleware.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Route("/api/v1", func(api chi.Router) {
		api.Use(render.SetContentType(render.ContentTypeJSON))

		api.Get("/status", func(w http.ResponseWriter, r *http.Request) {
			render.Status(r, http.StatusOK)
			render.JSON(w, r, struct {
				Status string `json:"status"`
			}{Status: "up and running"})
		})

		api.Get("/identity", s.Identity.GetIdentity)

		api.Route("/claims", func(rclaim chi.Router) {
			rclaim.Post("/", s.Claim.CreateClaim)
			rclaim.Get("/{id}", s.Claim.GetClaimByID)
			rclaim.Get("/revocation/status/{nonce}", s.Claim.RevocationStatus)
		})

		api.Route("/agent", func(internal chi.Router) {
			internal.Post("/", s.IDEN3Comm.Handle)
		})

	})

	return r
}
