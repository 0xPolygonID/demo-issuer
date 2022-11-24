package cfgs

type IssuerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL" yaml:"log_level"`

	DBFilePath string `mapstructure:"DB_FILE_PATH" yaml:"db_file_path"`
	ResetDb    bool   `mapstructure:"RESET_DB" yaml:"reset_db"`

	LocalUrl  string `mapstructure:"LOCAL_URL" yaml:"local_url"`
	PublicUrl string `mapstructure:"PUBLIC_URL" yaml:"public_url"`

	NodeRpcUrl                string `mapstructure:"NODE_RPC_URL" yaml:"node_rpc_url"`
	PublishingContractAddress string `mapstructure:"PUBLISHING_CONTRACT_ADDRESS" yaml:"publishing_contract_address"`
	PublishingPrivateKey      string `mapstructure:"PUBLISHING_PRIVATE_KEY" yaml:"publishing_private_key"`

	CircuitsDir       string `mapstructure:"CIRCUITS_DIR" yaml:"circuits_dir"`
	IpfsUrl           string `mapstructure:"IPFS_URL" yaml:"ipfs_url"`
	IdentitySecretKey string `mapstructure:"IDENTITY_SECRET_KEY" yaml:"identity_secret_key"`
}
