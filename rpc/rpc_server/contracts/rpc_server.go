package contracts

import (
	"context"
	"net"
	"time"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"
)

type RPCServer interface {
	WithOptions(options *options.RPCServerOptions) RPCServer
	WithHost(host string) RPCServer
	WithPort(port string) RPCServer
	WithServerShutdownTimeout(serverShutdownTimeout time.Duration) RPCServer
	WithServices(services []Services) RPCServer

	Start(listener net.Listener) error
	Shutdown(ctx context.Context) error
	Instance() any
	Address() string

	Logger() loggerContracts.Logger
}
