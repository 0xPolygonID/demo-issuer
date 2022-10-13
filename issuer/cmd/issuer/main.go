package main

import (
	"context"
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	app2 "lightissuer/issuer/app"
	"lightissuer/issuer/app/configs"
	handlers2 "lightissuer/issuer/app/handlers"
	services2 "lightissuer/issuer/services"
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

	schemaService := services2.NewSchemaService(config.URL)
	dbService, err := services2.NewDBService()

	if err != nil {
		log.Fatal("failed starting DB service", err)
		return
	}

	// will automatically create an Identity as well for the issuer
	iService, err := services2.NewIdentityService(ctx, &privKey, dbService, schemaService)
	if err != nil {
		log.Fatal("error starting identity service ", err)
	}

	claimService := services2.NewClaimService(iService, dbService, *schemaService)
	if err != nil {
		log.Fatal("couldn't get verification key")
		return
	}

	h := app2.Handlers{Identity: handlers2.NewIdentityHandler(iService), Claim: handlers2.NewClaimHandler(dbService, claimService, schemaService, iService, config),
		IDEN3Comm: handlers2.NewIDEN3CommHandler(iService, dbService, claimService, schemaService, config.CircuitPath),
	}

	router := h.Routes()
	server := app2.NewServer(router)

	fmt.Printf("demo-issuer started on %s \n", config.Host)

	err = server.Run(config.Port)
	if err != nil {
		log.Fatal("error in starting server")
	}

}
