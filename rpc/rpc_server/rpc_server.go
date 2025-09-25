package rpcserver

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
	"frisboo-bank/pkg/validation"
)

var _ contracts.RPCServer = (*rpcServer)(nil)

type rpcServer struct {
	adapter contracts.RPCServerAdapter
}

func New(adapter contracts.RPCServerAdapter) contracts.RPCServer {
	validation.Assert(adapter != nil, "adapter can't be nil")

	return &rpcServer{adapter}
}

func (r *rpcServer) Stop(ctx context.Context) error {
	return r.adapter.Stop(ctx)
}

func (r *rpcServer) Start(ctx context.Context) error {
	return r.adapter.Start(ctx)
}

func (r *rpcServer) Name() string {
	return r.adapter.Name()
}

func (r *rpcServer) Type() rpcservertype.RpcServerType {
	return r.adapter.Type()
}

func (r *rpcServer) Config() *config.Config {
	return r.adapter.Config()
}

func (r *rpcServer) Logger() loggerContracts.Logger {
	return r.adapter.Logger()
}
