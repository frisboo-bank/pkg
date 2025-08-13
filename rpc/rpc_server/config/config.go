package config

import (
	"net"
	"strings"
	"time"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

var pError = customerrors.PrefixedError("rpc-server config")

type EnvConfig struct {
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	ServerShutdownTimeout time.Duration               `mapstructure:"serverShutdownTimeout"`
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
}

func LoadEnvConfig(loader configContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.LoadConfig[EnvConfig](loader, env, "rpcServer")
}

type Config struct {
	Host                  string
	Port                  string
	ServerShutdownTimeout time.Duration
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

var defaultConfig = &Config{
	Host:                  "0.0.0.0",
	Port:                  "9000",
	ServerShutdownTimeout: constants.SERVER_SHUTDOWN_TIMEOUT,
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(&Config{})
}

func FromEnvConfig(cfg *EnvConfig) *options.OptionBuilder[Config] {
	opts := Apply()

	if cfg.Host != "" {
		opts.With(Host(cfg.Host))
	}

	if cfg.Port != "" {
		opts.With(Port(cfg.Port))
	}

	if cfg.ServerShutdownTimeout != 0 {
		opts.With(ServerShutdownTimeout(cfg.ServerShutdownTimeout))
	}

	return opts
}

func Host(host string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		host = strings.TrimSpace(host)

		if host == "" {
			return pError.New("host cannot be empty")
		}

		cfg.Host = host
		return nil
	})
}

func Port(port string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		port = strings.TrimSpace(port)

		if port == "" {
			return pError.New("port cannot be empty")
		}

		cfg.Port = port
		return nil
	})
}

func ServerShutdownTimeout(serverShutdownTimeout time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if serverShutdownTimeout <= 0 {
			return pError.New(" serverShutdownTimeout must be positive")
		}

		cfg.ServerShutdownTimeout = serverShutdownTimeout
		return nil
	})
}
