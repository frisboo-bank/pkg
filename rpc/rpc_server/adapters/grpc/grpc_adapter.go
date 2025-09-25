package grpc

import (
	"context"
	"errors"
	"net"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
	"frisboo-bank/pkg/validation"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	googlerpc "google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

var _ contracts.RPCServerAdapter = (*grpcRPCServerAdapter)(nil)

type grpcRPCServerAdapter struct {
	name     string
	cfg      *config.Config
	listener net.Listener
	logger   loggerContracts.Logger
	server   *googlerpc.Server
}

func New(name string, cfg *config.Config, logger loggerContracts.Logger, metrics any) contracts.RPCServerAdapter {
	validation.AssertNotNil("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	server := googlerpc.NewServer(
		googlerpc.StatsHandler(otelgrpc.NewServerHandler()),
		// googlerpc.StatsHandler(otel.NewServerHandler()),

		googlerpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     0,
			MaxConnectionAge:      0,
			MaxConnectionAgeGrace: 0,
			Time:                  0,
			Timeout:               0,
		}),

		googlerpc.StreamInterceptor(grpcMiddleware.ChainStreamServer()),

		googlerpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			grpcCtxTags.UnaryServerInterceptor(),
			grpcRecovery.UnaryServerInterceptor(),
		)),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("test", grpc_health_v1.HealthCheckResponse_SERVING)
	// reflection.Register(server)

	return &grpcRPCServerAdapter{
		name:     name,
		cfg:      cfg,
		listener: nil,
		server:   server,
		logger:   logger,
	}
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

	reflection.Register(g.server)

	if err := g.server.Serve(listener); err != nil && !errors.Is(err, googlerpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (g *grpcRPCServerAdapter) Stop(ctx context.Context) error {
	g.server.GracefulStop()
	return nil
}

func (g *grpcRPCServerAdapter) Name() string {
	return g.name
}

func (g *grpcRPCServerAdapter) Type() rpcservertype.RpcServerType {
	return rpcservertype.RpcServerTypes.GRPC
}

func (g *grpcRPCServerAdapter) Config() *config.Config {
	return g.cfg
}

func (g *grpcRPCServerAdapter) Logger() loggerContracts.Logger {
	return g.logger
}
