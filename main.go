package main

import (
	"context"
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"lightissuer/app"
	"lightissuer/app/configs"
	"lightissuer/app/handlers"
	"lightissuer/services"
	"log"
)

func main() {

	ctx := context.Background()
	config, err := configs.ParseIssuerConfig()
	if err != nil {
		log.Fatal("error in parsing config")
	}

	var privKey babyjub.PrivateKey
	copy(privKey[:], config.PrivateKey)

	schemaService := services.NewSchemaService(config.URL)
	dbService, err := services.NewDBService()

	if err != nil {
		log.Fatal("failed starting DB service", err)
		return
	}

	// will automatically create an Identity as well for the issuer
	iService, err := services.NewIdentityService(ctx, &privKey, dbService, schemaService)
	if err != nil {
		log.Fatal("error starting identity service ", err)
	}

	claimService := services.NewClaimService(iService, dbService, *schemaService)
	if err != nil {
		log.Fatal("couldn't get verification key")
		return
	}

	h := app.Handlers{Identity: handlers.NewIdentityHandler(iService), Claim: handlers.NewClaimHandler(dbService, claimService, schemaService, iService, config),
		IDEN3Comm: handlers.NewIDEN3CommHandler(iService, dbService, claimService, schemaService, config.CircuitPath),
	}

	router := h.Routes()
	server := app.NewServer(router)

	fmt.Printf("demo-issuer started on %s \n", config.Host)

	err = server.Run(config.Port)
	if err != nil {
		log.Fatal("error in starting server")
	}

}
