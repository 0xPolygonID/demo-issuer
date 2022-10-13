package main

import (
	"flag"
	"fmt"
	"issuer/service"
	"os"
)

func main() {
	var fileName string
	flag.StringVar(&fileName, "cfg-file", "", "an alternative path to the cfg file")
	flag.Parse()

	if err := run(fileName); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run(fileName string) error {

	return service.CreateApp(fileName)
}

//func main() {
//
//	ctx := context.Background()
//	cfgs, err := configs.ParseIssuerConfig()
//	if err != nil {
//		log.Fatal("error in parsing cfgs")
//	}
//
//	var privKey babyjub.PrivateKey	// put inside cfg
//	copy(privKey[:], cfgs.PrivateKey)
//
//	schemaService := services2.NewSchemaService(cfgs.URL)
//	dbService, err := services2.NewDBService()
//
//	if err != nil {
//		log.Fatal("failed starting DB service", err)
//		return
//	}
//
//	// will automatically create an Identity as well for the issuer
//	iService, err := services2.NewIdentityService(ctx, &privKey, dbService, schemaService)
//	if err != nil {
//		log.Fatal("error starting identity service ", err)
//	}
//
//	claimService := services2.NewClaimService(iService, dbService, *schemaService)
//	if err != nil {
//		log.Fatal("couldn't get verification key")
//		return
//	}
//
//	h := app2.Handlers{Identity: handlers2.NewIdentityHandler(iService), Claim: handlers2.NewClaimHandler(dbService, claimService, schemaService, iService, cfgs),
//		IDEN3Comm: handlers2.NewIDEN3CommHandler(iService, dbService, claimService, schemaService, cfgs.CircuitPath),
//	}
//
//	router := h.Routes()
//	server := app2.NewServer(router)
//
//	fmt.Printf("demo-issuer started on %s \n", cfgs.Host)
//
//	err = server.Run(cfgs.Port)
//	if err != nil {
//		log.Fatal("error in starting server")
//	}
//
//}
