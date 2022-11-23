package cfgs

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func New(altCfgPath string) (*IssuerConfig, error) {
	cfgFilePath, err := resolveConfigPath(altCfgPath)
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(cfgFilePath); !os.IsNotExist(err) {
		viper.SetConfigFile(cfgFilePath)
	}

	setViperEnvConfig()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Infof("using config file: %s", viper.ConfigFileUsed())
	}

	cfgObj := &IssuerConfig{}
	if err := viper.Unmarshal(cfgObj); err != nil {
		return nil, err
	}

	logger.Infof("issuer config: %v", cfgObj)
	return cfgObj, nil
}

// resolveConfigPath resolves the configuration
// file path (YAML formatted)
// It's looking for asset-communication-config.yaml in the current
// directory (higher precedence) and then in $HOME/.qed-it (lower precedence)
func resolveConfigPath(cfgFilePath string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(cfgFilePath) == 0 {
		cfgFilename := "issuer_config.default.yaml"
		cwdCfgPath := fmt.Sprintf("%s/%s", cwd, cfgFilename)

		_, err = os.Stat(cwdCfgPath)

		return cwdCfgPath, err
	}

	cwdCfgPath := fmt.Sprintf("%s/%s", cwd, cfgFilePath)
	_, err = os.Stat(cwdCfgPath)

	return cwdCfgPath, nil
}

func setViperEnvConfig() {
	viper.SetEnvPrefix("ISSUER_SERVICE")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	setConfigDefaults()
	viper.AutomaticEnv() // read in environment variables that match
}

func setConfigDefaults() {
	viper.SetDefault("LOG_LEVEL", "INFO")
	viper.SetDefault("DB_FILE_PATH", "issuer.db")
	viper.SetDefault("KEY_DIR", "keys")
	viper.SetDefault("LOCAL_URL", "localhost:8001")
	viper.SetDefault("IPFS_URL", "ipfs.io")
}
