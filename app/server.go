package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"net/http"
)

// Server app server
type Server struct {
	Routes     chi.Router
	httpServer *http.Server
}

// NewServer create new server
func NewServer(router chi.Router) *Server {
	return &Server{
		Routes: router,
	}
}

// Close closes connection to the server
func (s *Server) Close(ctx context.Context) error {
	return errors.WithStack(s.httpServer.Shutdown(ctx))
}

// Run server
func (s *Server) Run(port int) error {
	println("Server listening on port", port)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.Routes}
	return errors.WithStack(s.httpServer.ListenAndServe())
}
