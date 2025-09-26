package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
)

type (
	rpcServerCommon interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		Name() string
		Type() rpcservertype.RpcServerType
		Config() *config.Config
		Logger() loggerContracts.Logger
	}

	RPCServer interface {
		rpcServerCommon
	}

	RPCServerAdapter interface {
		rpcServerCommon
	}
)
