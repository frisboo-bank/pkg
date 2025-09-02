package config

import (
	"net"
	"time"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"

	loggerConfig "frisboo-bank/pkg/logger/config"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	grpcConfig "frisboo-bank/pkg/rpc/rpc_server/adapters/grpc/config"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	ServerShutdownTimeout time.Duration               `mapstructure:"serverShutdownTimeout"`

	// adapters
	GRPC *grpcConfig.Config `mapstructure:"grpc"`

	// dependency
	Logger *loggerConfig.Config `mapstructure:"logger"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() *Config {
	return &Config{
		Host:                  "0.0.0.0",
		Port:                  "9000",
		ServerShutdownTimeout: 30 * time.Second,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
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
