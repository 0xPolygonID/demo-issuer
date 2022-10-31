package main

import (
	"fmt"
	"issuer/backend/handler"
	"issuer/cfgs"
	"log"
	"net/http"
)

func main() {
	cfg, err := cfgs.New("")
	if err != nil {
		log.Fatal("error loading cfg file")
	}

	h := handler.NewHandler(cfg)
	mux := http.NewServeMux()
	mux.HandleFunc("/api/sign-in", h.GetQR)
	mux.HandleFunc("/api/callback", h.Callback)
	mux.HandleFunc("/api/status", h.Status)

	// TODO: should be configured from environment
	port := "3000"

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	fmt.Println("Running backend service on 3000")

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen:%+s\n", err)
	}
}
