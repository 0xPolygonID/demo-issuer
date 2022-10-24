package cfgs

type IssuerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL" yaml:"log_level"`

	HttpListenAddress string   `mapstructure:"HTTP_LISTEN_ADDRESS" yaml:"http_listen_address"`
	DB                DBConfig `mapstructure:"DB" yaml:"db"`
	MerkleTreeDepth   int      `mapstructure:"MERKLE_TREE_DEPTH" yaml:"merkle_tree_depth"`
	SecretKey         []byte   `mapstructure:"SECRET_KEY" yaml:"secret_key"`
}

type DBConfig struct {
	Kind     string `mapstructure:"KIND" yaml:"kind"`
	FilePath string `mapstructure:"FILE_PATH" yaml:"file_path"`
}
