package configs

import "github.com/kelseyhightower/envconfig"

type Environement struct {
	Port        int    `envconfig:"PORT" default:"8001"`
	Host        string `envconfig:"HOST" default:""`
	URL         string `envconfig:"IPFS_URL" default:"ipfs.io"`
	CircuitPath string `envconfig:"CIRCUIT_PATH" default:"circuits"`
	PrivateKey  []byte `envconfig:"PRIVATE_KEY" default:"48,35,125,91,136,34,33,104,95,41,38,54,28,167,232,31,18,53,213,250,239,2,188,251,249,159,69,151,24,191,217,82"`
}

// ParseIssuerConfig parses env variables in to Issuer config.
func ParseIssuerConfig() (*Environement, error) {
	var cfg Environement
	err := envconfig.Process("ISSUER", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
