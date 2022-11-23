package service

import (
	"encoding/hex"
	"github.com/iden3/go-iden3-crypto/babyjub"
	logger "github.com/sirupsen/logrus"
	database "issuer/db"
	"issuer/service/blockchain"
	"issuer/service/cfgs"
	"issuer/service/http"
	"issuer/service/identity"
	"issuer/service/identity/state"
	"issuer/service/schema"
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
	idenState, err := state.NewIdentityState(db)
	if err != nil {
		return err
	}

	logger.Info("processing secret key")
	sk, err := secretKeyToBabyJub(cfg.Identity.SecretKey)
	if err != nil {
		return err
	}

	schemaBuilder := schema.NewBuilder(cfg.IpfsUrl)

	chstore, err := blockchain.NewBlockchainConnect(cfg.NodeRpcUrl, cfg.ContractAddress, cfg.Identity.PublishingKey)
	if err != nil {
		return err
	}

	logger.Info("creating Identity")
	issuer, err := identity.New(idenState, schemaBuilder, sk, cfg, chstore)
	if err != nil {
		return err
	}

	s := http.NewServer(cfg.LocalUrl, issuer)

	logger.Infof("spining up API server @%s", cfg.LocalUrl)
	return s.Run()
}

func secretKeyToBabyJub(sk string) (babyjub.PrivateKey, error) {
	if len(sk) == 0 {
		return babyjub.NewRandPrivKey(), nil
	}

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
