package service

import (
	"encoding/hex"
	"github.com/iden3/go-iden3-crypto/babyjub"
	logger "github.com/sirupsen/logrus"
	database "issuer/db"
	"issuer/service/backend"
	"issuer/service/cfgs"
	"issuer/service/command"
	"issuer/service/http"
	"issuer/service/identity"
	"issuer/service/identity/state"
	"os"
)

func CreateApp(altCfgPath string) error {
	logger.Info("boot up issuer service")

	logger.Info("loading Configuration")
	cfg, err := cfgs.New(altCfgPath)
	if err != nil {
		return err
	}

	logger.Info("setting up logger")
	err = initGlobalLogger(cfg.LogLevel)
	if err != nil {
		return err
	}

	logger.Info("creating DB")
	db, err := database.New(cfg.DBFilePath)
	if err != nil {
		return err
	}

	logger.Info("creating identity state")
	idenState, err := state.NewIdentityState(db, cfg.Identity.MerkleTreeDepth)
	if err != nil {
		return err
	}

	logger.Info("processing secret key")
	sk, err := secretKeyToBabyJub(cfg.Identity.SecretKey)
	if err != nil {
		return err
	}

	logger.Info("creating Identity")
	issuer, err := identity.New(idenState, sk, cfg.IssuerUrl)
	if err != nil {
		return err
	}

	commHandler := command.NewHandler(idenState, cfg.CircuitPath)
	backendHandler := backend.NewHandler(cfg)

	// start service
	s := http.NewServer(cfg.Http.Host, cfg.Http.Port, issuer, commHandler, backendHandler)

	logger.Infof("spining up API server @%s", cfg.Http.Host+":"+cfg.Http.Port)
	return s.Run()
}

func secretKeyToBabyJub(sk string) (babyjub.PrivateKey, error) {
	b, err := hex.DecodeString(sk)
	if err != nil {
		return babyjub.PrivateKey{}, err
	}
	var privKey babyjub.PrivateKey
	copy(privKey[:], b)

	return privKey, nil
}

func initGlobalLogger(level string) error {
	logLevel, err := logger.ParseLevel(level)
	if err != nil {
		return err
	}

	logger.SetLevel(logLevel)
	//logger.SetFormatter(&logger.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetOutput(os.Stdout)
	return nil
}
