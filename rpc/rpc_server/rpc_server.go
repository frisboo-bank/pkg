package rpcserver

import (
	"context"

	"frisboo-bank/pkg/customerrors"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
	"frisboo-bank/pkg/utils"
)

var _ contracts.RPCServer = (*rpcServer)(nil)

var pError = customerrors.PrefixedError("rpc server")

type rpcServer struct {
	cfg     *config.Config
	adapter contracts.RPCServerAdapter
	logger  loggerContracts.Logger
}

func New(
	adapter contracts.RPCServerAdapter,
	logger loggerContracts.Logger,
	opts *options.OptionBuilder[config.Config],
) (contracts.RPCServer, error) {
	utils.Assert(adapter != nil, pError.New("adapter can't be nil"))
	utils.Assert(logger != nil, pError.New("logger can't be nil"))
	utils.Assert(opts != nil, pError.New("opts can't be nil"))

	cfg := opts.Build()

	server := &rpcServer{
		cfg:     cfg,
		adapter: adapter,
		logger:  logger,
	}

	if err := adapter.Setup(cfg); err != nil {
		return nil, err
	}

	return server, nil
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
	return r.logger
}
