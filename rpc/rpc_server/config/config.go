package config

import (
	"time"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"

	configContracts "frisboo-bank/pkg/config/contracts"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

var (
	Type                  = rpcservertype.RpcServerTypes.GRPC
	Host                  = "0.0.0.0"
	Port                  = "9000"
	ServerShutdownTimeout = constants.SERVER_SHUTDOWN_TIMEOUT
)

type RPCServerConfig struct {
	Type                  rpcservertype.RpcServerType `mapstructure:"type"`
	Host                  string                      `mapstructure:"host"`
	Port                  string                      `mapstructure:"port"`
	ServerShutdownTimeout time.Duration               `mapstructure:"serverShutdownTimeout"`
}

func ProvideRPCServerConfig(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*RPCServerConfig, error) {
	return config.LoadConfig[RPCServerConfig](loader, env, "rpcServer")
}
