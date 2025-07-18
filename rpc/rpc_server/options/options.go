package options

import (
	"time"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
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
	Services              contracts.Services
	Logger                loggerContracts.Logger
}

func ProvideRPCServerOptions(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*RPCServerOptions, error) {
	return config.LoadOptions[RPCServerOptions](loader, env)
}
