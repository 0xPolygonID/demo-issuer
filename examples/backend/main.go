package main

import (
	"auth_flow/handlers"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/sign-in", handlers.GetQR)
	mux.HandleFunc("/api/callback", handlers.Callback)
	mux.HandleFunc("/api/status", handlers.Status)

	// TODO: should be configured from environment
	port := "3000"

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen:%+s\n", err)
	}
}
