package service

import (
	logger "github.com/sirupsen/logrus"
	"issuer/service/cfgs"
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

	// create main components

	// init an api server

	return nil
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
