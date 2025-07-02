package contracts

import (
	"context"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	"net"
)

type RPCServer interface {
	Start(listener net.Listener) error
	Shutdown(ctx context.Context) error

	Logger() loggerContracts.Logger
	Config() *options.RPCServerOptions
}
