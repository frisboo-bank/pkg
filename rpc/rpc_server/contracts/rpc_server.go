package contracts

import (
	"context"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	"net"
	"time"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
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
