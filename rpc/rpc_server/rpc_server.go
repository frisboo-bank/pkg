package rpcserver

import (
	"context"

	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/syserrors"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

var _ contracts.RPCServer = (*rpcServer)(nil)

type rpcServer struct {
	adapter contracts.RPCServerAdapter
}

func New(adapter contracts.RPCServerAdapter) contracts.RPCServer {
	syserrors.Assert(adapter != nil, "adapter can't be nil")

	return &rpcServer{adapter}
}

func (r *rpcServer) Shutdown(ctx context.Context) error {
	return r.adapter.Shutdown(ctx)
}

func (r *rpcServer) Start(ctx context.Context) error {
	return r.adapter.Start(ctx)
}

func (r *rpcServer) Type() rpcservertype.RpcServerType {
	return r.adapter.Type()
}

func (r *rpcServer) Logger() loggerContracts.Logger {
	return r.adapter.Logger()
}
