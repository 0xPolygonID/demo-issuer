package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

type Server struct {
	Routes     chi.Router
	httpServer *http.Server
}

func NewServer(router chi.Router) *Server {
	return &Server{
		Routes: router,
	}
}

func (s *Server) Close(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Run(port int) error {
	println("Server listening on port", port)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.Routes}
	return s.httpServer.ListenAndServe()
}
