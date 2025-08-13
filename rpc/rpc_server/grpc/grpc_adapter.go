package grpc

import (
	"context"
	"errors"
	"net"

	"frisboo-bank/pkg/customerrors"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
	"frisboo-bank/pkg/utils"

	googlerpc "google.golang.org/grpc"
)

var _ contracts.RPCServerAdapter = (*grpcRPCServerAdapter)(nil)

var pError = customerrors.PrefixedError("grpc server")

type grpcRPCServerAdapter struct {
	cfg      *config.Config
	listener net.Listener
	logger   loggerContracts.Logger
	server   *googlerpc.Server
}

func New(logger loggerContracts.Logger) contracts.RPCServerAdapter {
	utils.Assert(logger != nil, pError.New("logger can't be nil"))

	return &grpcRPCServerAdapter{
		logger: logger,
	}
}

func (g *grpcRPCServerAdapter) Setup(cfg *config.Config) error {
	g.cfg = cfg
	g.server = googlerpc.NewServer()

	return nil
}

func (g *grpcRPCServerAdapter) Start(ctx context.Context) error {
	g.logger.Info("starting server...")

	addr := g.cfg.Address()

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g.listener = listener

	g.logger.Infof("server listening on address: %s", addr)

	if err := g.server.Serve(listener); err != nil && !errors.Is(err, googlerpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (g *grpcRPCServerAdapter) Shutdown(ctx context.Context) error {
	g.logger.Info("server shutting down...")

	g.server.GracefulStop()

	g.logger.Info("server shutdown done successfully")

	return nil
}

func (g *grpcRPCServerAdapter) Type() rpcservertype.RpcServerType {
	return rpcservertype.RpcServerTypes.GRPC
}
