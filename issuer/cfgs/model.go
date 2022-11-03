package cfgs
package cfgs
package cfgs

type IssuerConfig struct {
	LogLevel    string `mapstructure:"LOG_LEVEL" yaml:"log_level"`
	DBFilePath  string `mapstructure:"DB_FILE_PATH" yaml:"db_file_path"`
	CircuitPath string `mapstructure:"CIRCUIT_PATH" yaml:"circuit_path"`
	KeyDir      string `mapstructure:"KEY_DIR" yaml:"key_dir"`
	HostUrl     string `mapstructure:"HOST_URL" yaml:"host_url"`
	IssuerUrl   string `mapstructure:"ISSUER_URL" yaml:"issuer_url"`

	Http     Http     `mapstructure:"HTTP" yaml:"http"`
	Identity Identity `mapstructure:"IDENTITY" yaml:"identity"`
}

type Http struct {
	Host string `mapstructure:"HOST" yaml:"host"`
	Port string `mapstructure:"PORT" yaml:"port"`
}

type Identity struct {
	SecretKey       string `mapstructure:"SECRET_KEY" yaml:"secret_key"`
	MerkleTreeDepth int    `mapstructure:"MERKLE_TREE_DEPTH" yaml:"merkle_tree_depth"`
	//IdentityHostUrl string `mapstructure:"IDENTITY_HOST_URL" yaml:"identity_host_url"`
}
