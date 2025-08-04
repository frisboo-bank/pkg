package contracts

import (
	"context"
	"net"
	"time"

	"frisboo-bank/pkg/rpc/rpc_server/config"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

type (
	rpcServerConfig interface {
		WithConfig(config *config.RPCServerConfig) RPCServer
		WithHost(host string) RPCServer
		WithPort(port string) RPCServer
		WithServerShutdownTimeout(serverShutdownTimeout time.Duration) RPCServer
		WithServices(services []Services) RPCServer
		WithLogger(logger loggerContracts.Logger) RPCServer

		Host() string
		Port() string
		ServerShutdownTimeout() time.Duration
		Services() []Services
		Logger() loggerContracts.Logger
	}

	rpcServerCore interface {
		Start(listener net.Listener) error
		Shutdown(ctx context.Context) error
		Address() string
		Instance() any
		Type() rpcservertype.RpcServerType
	}

	RPCServer interface {
		rpcServerConfig
		rpcServerCore
	}

	RPCServerInternal interface {
		SetupInstance()
	}
)
