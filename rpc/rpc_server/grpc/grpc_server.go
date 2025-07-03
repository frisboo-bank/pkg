package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	"frisboo-bank/pkg/utils"

	googlerpc "google.golang.org/grpc"
)

type GRPCServer struct {
	grpcServer *googlerpc.Server
	logger     loggerContracts.Logger
	config     *options.RPCServerOptions
}

var _ contracts.RPCServer = (*GRPCServer)(nil)

func NewGRPCServer(config *options.RPCServerOptions) contracts.RPCServer {
	utils.Assert(config.Logger != nil, "(rpc-server) logger must not be nil")

	return newGRPCServer(config)
}

func newGRPCServer(config *options.RPCServerOptions) contracts.RPCServer {
	return &GRPCServer{
		grpcServer: googlerpc.NewServer(),
		logger:     config.Logger,
		config:     config,
	}
}

func (g *GRPCServer) Start(listener net.Listener) error {
	if err := g.grpcServer.Serve(listener); err != nil && !errors.Is(err, googlerpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (g *GRPCServer) Shutdown(ctx context.Context) error {
	if g.grpcServer == nil {
		return fmt.Errorf("grpc-server: looks like there is no server running")
	}

	g.grpcServer.GracefulStop()

	return nil
}

func (g *GRPCServer) Config() *options.RPCServerOptions {
	return g.config
}

func (g *GRPCServer) Logger() loggerContracts.Logger {
	return g.logger
}
