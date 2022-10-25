package cfgs

type IssuerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL" yaml:"log_level"`

	HttpListenAddress string `mapstructure:"HTTP_LISTEN_ADDRESS" yaml:"http_listen_address"`
	IdentityHostUrl   string `mapstructure:"IDENTITY_HOST_URL" yaml:"identity_host_url"`
	DBFilePath        string `mapstructure:"DB_FILE_PATH" yaml:"db_file_path"`
	MerkleTreeDepth   int    `mapstructure:"MERKLE_TREE_DEPTH" yaml:"merkle_tree_depth"`
	SecretKey         []byte `mapstructure:"SECRET_KEY" yaml:"secret_key"`
}
