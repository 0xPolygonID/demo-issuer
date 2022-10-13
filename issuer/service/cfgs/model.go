package cfgs

type IssuerConfig struct {
	LogLevel string `mapstructure:"LOG_LEVEL" yaml:"log_level"`

	HttpListenAddress string   `mapstructure:"HTTP_LISTEN_ADDRESS" yaml:"http_listen_address"`
	DB                DBConfig `mapstructure:"DB" yaml:"db"`
}

type DBConfig struct {
	Kind             string `mapstructure:"KIND" yaml:"kind"`
	ConnectionString string `mapstructure:"CONNECTION_STRING" yaml:"connection_string"`
}