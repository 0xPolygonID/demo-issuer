package cfgs

type IssuerConfig struct {
	LogLevel   string   `mapstructure:"LOG_LEVEL" yaml:"log_level"`
	DBFilePath string   `mapstructure:"DB_FILE_PATH" yaml:"db_file_path"`
	KeyDir     string   `mapstructure:"KEY_DIR" yaml:"key_dir"`
	PublicUrl  string   `mapstructure:"PUBLIC_URL" yaml:"public_url"`
	LocalUrl   string   `mapstructure:"LOCAL_URL" yaml:"local_url"`
	Identity   Identity `mapstructure:"IDENTITY" yaml:"identity"`
	IpfsUrl    string   `mapstructure:"IPFS_URL" yaml:"ipfs_url"`
	NodeRpcUrl string   `mapstructure:"NODE_RPC_URL" yaml:"node_rpc_url"`
}

type Identity struct {
	SecretKey string `mapstructure:"SECRET_KEY" yaml:"secret_key"`
}
