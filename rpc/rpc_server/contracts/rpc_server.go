package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
)

type (
	RPCServer interface {
		Start(ctx context.Context) error
		Shutdown(ctx context.Context) error
		Type() rpcservertype.RpcServerType
		Logger() loggerContracts.Logger
	}

	RPCServerAdapter interface {
		RPCServer
	}
)
