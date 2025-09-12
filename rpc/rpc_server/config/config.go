package config

import (
	"net"
	"time"

	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"

	cValidation "frisboo-bank/pkg/validation"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	grpcConfig "frisboo-bank/pkg/rpc/rpc_server/adapters/grpc/config"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationIs "github.com/go-ozzo/ozzo-validation/v4/is"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	ServerShutdownTimeout time.Duration               `mapstructure:"serverShutdownTimeout"`

	// adapters
	GRPC *grpcConfig.Config `mapstructure:"grpc"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c *Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() Config {
	return Config{
		Type:                  rpcservertype.RpcServerTypes.GRPC,
		Host:                  "0.0.0.0",
		Port:                  "9000",
		ServerShutdownTimeout: 30 * time.Second,

		// adapters
		GRPC: grpcConfig.Default(),
	}
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(rpcservertype.RpcServerTypes))),
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.ServerShutdownTimeout, validation.Required, validation.Min(0)),
	); err != nil {
		return err
	}

	switch c.Type {
	case rpcservertype.RpcServerTypes.GRPC:
		if err := validation.Validate(&c.GRPC, validation.Required); err != nil {
			return err
		}
		return c.GRPC.Validate()
	}

	return nil
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"rpcServers",
		"rpcServer",
		Default,
	)
}
