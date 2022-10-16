package service

import (
	logger "github.com/sirupsen/logrus"
	"issuer/service/cfgs"
	"issuer/service/claims"
	database "issuer/service/db"
	"issuer/service/http"
	"issuer/service/identity"
	"os"
	"time"
)

func CreateApp(altCfgPath string) error {

	cfg, err := cfgs.New(altCfgPath)
	if err != nil {
		return err
	}

	// setup log
	err = initGlobalLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	// create dependencies

	db, err := database.New(cfg.DB.FilePath)
	if err != nil {
		return err
	}

	// create main components
	// identity
	// claims
	// agent

	claimsManager, err := claims.New()
	if err != nil {
		return err
	}

	identityManager, err := identity.New(db, )
	if err != nil {
		return err
	}

	//agent, err := agent

	// init an api server
	s := http.NewServer(cfg.HttpListenAddress, identityManager, claimsManager)
	return s.Run()
}

func initGlobalLogger(level string) error {
	logLevel, err := logger.ParseLevel(level)
	if err != nil {
		return err
	}

	logger.SetLevel(logLevel)
	logger.SetFormatter(&logger.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetOutput(os.Stdout)
	return nil
}
