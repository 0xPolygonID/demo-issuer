package cfgs

type IssuerConfig struct {
	LogLevel        string   `mapstructure:"LOG_LEVEL" yaml:"log_level"`
	DBFilePath      string   `mapstructure:"DB_FILE_PATH" yaml:"db_file_path"`
	KeyDir          string   `mapstructure:"KEY_DIR" yaml:"key_dir"`
	PublicUrl       string   `mapstructure:"PUBLIC_URL" yaml:"public_url"`
	LocalUrl        string   `mapstructure:"LOCAL_URL" yaml:"local_url"`
	Identity        Identity `mapstructure:"IDENTITY" yaml:"identity"`
	IpfsUrl         string   `mapstructure:"IPFS_URL" yaml:"ipfs_url"`
	ContractAddress string   `mapstructure:"CONTRACT_ADDRESS" yaml:"contract_address"`
	NodeURL         string   `mapstructure:"NODE_URL" yaml:"node_url"`
}

type Identity struct {
	SecretKey       string `mapstructure:"SECRET_KEY" yaml:"secret_key"`
	MerkleTreeDepth int    `mapstructure:"MERKLE_TREE_DEPTH" yaml:"merkle_tree_depth"`
}
