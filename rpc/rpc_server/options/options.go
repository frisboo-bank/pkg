package options

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"time"

	configContracts "frisboo-bank/pkg/config/contracts"
)

var (
	Host                  = "0.0.0.0"
	Port                  = "9000"
	ServerShutdownTimeout = constants.SERVER_SHUTDOWN_TIMEOUT
)

type RPCServerOptions struct {
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout"`
}

func ProvideRPCServerOptions(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*RPCServerOptions, error) {
	return config.LoadOptions[RPCServerOptions](loader, env)
}
