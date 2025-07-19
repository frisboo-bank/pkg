package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"

	googlerpc "google.golang.org/grpc"
)

type GRPCServer struct {
	host                  string
	port                  string
	serverShutdownTimeout time.Duration
	services              []contracts.Services

	grpcServer *googlerpc.Server
	logger     loggerContracts.Logger
}

func (g *GRPCServer) WithOptions(options *options.RPCServerOptions) contracts.RPCServer {
	return g.
		WithHost(options.Host).
		WithPort(options.Port).
		WithServerShutdownTimeout(options.ServerShutdownTimeout)
}

func (g *GRPCServer) WithHost(host string) contracts.RPCServer {
	g.host = host
	return g
}

func (g *GRPCServer) WithPort(port string) contracts.RPCServer {
	g.port = port
	return g
}

func (g *GRPCServer) WithServerShutdownTimeout(serverShutdownTimeout time.Duration) contracts.RPCServer {
	g.serverShutdownTimeout = serverShutdownTimeout
	return g
}

func (g *GRPCServer) WithServices(services []contracts.Services) contracts.RPCServer {
	g.services = services
	return g
}

var _ contracts.RPCServer = (*GRPCServer)(nil)

func NewGRPCServer(logger loggerContracts.Logger) contracts.RPCServer {
	return &GRPCServer{
		host:                  options.Host,
		port:                  options.Port,
		serverShutdownTimeout: options.ServerShutdownTimeout,
		grpcServer:            googlerpc.NewServer(),
		logger:                logger,
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

func (g *GRPCServer) Instance() any {
	return g.grpcServer
}

func (g *GRPCServer) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

func (g *GRPCServer) Logger() loggerContracts.Logger {
	return g.logger
}
