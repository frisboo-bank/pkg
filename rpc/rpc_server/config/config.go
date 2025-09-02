package config

import (
	"net"
	"time"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"

	loggerConfig "frisboo-bank/pkg/logger/config"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

type Config struct {
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	ServerShutdownTimeout time.Duration               `mapstructure:"serverShutdownTimeout"`
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
	Logger                loggerConfig.Config         `mapstructure:"logger"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() *Config {
	loggerCfg := loggerConfig.Default()
	loggerCfg.Prefix = "rpc-server"

	return &Config{
		Host:                  "0.0.0.0",
		Port:                  "9000",
		ServerShutdownTimeout: 30 * time.Second,
		Logger:                *loggerCfg,
	}
}

var defaultConfig = &Config{
	Host:                  "0.0.0.0",
	Port:                  "9000",
	ServerShutdownTimeout: constants.SERVER_SHUTDOWN_TIMEOUT,
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("rpc-server", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
