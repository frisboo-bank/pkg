package rpcserver

import (
	"context"
	"errors"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/factory"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
	"net"

	"google.golang.org/grpc"
)

var Module = container.NewModule("rpc-server",
	container.Provide(func(env environment.Environment) (*options.RPCServerOptions, error) {
		return options.ProvideRPCServerOptions(env)
	}),
	container.Provide(func(config *options.RPCServerOptions, logger loggerContracts.Logger) (contracts.RPCServer, error) {
		logger = logger.WithName("rpc-server")

		return factory.GetInstance(config,
			options.WithLogger(logger),
		)
	}),
	container.Hook(startHook, stopHook),
)

func startHook(rpcServer contracts.RPCServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		addr := rpcServer.Config().Address()

		rpcServer.Logger().Info("(rpc-server) starting server...")

		go func() {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				rpcServer.Logger().Fatalf("(rpc-server) failed to listen on address with error: %v", err)
				return
			}

			if err := rpcServer.Start(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				rpcServer.Logger().Fatalf("(rpc-server) failed to start with error: %v", err)
			}
		}()

		rpcServer.Logger().Infof("(rpc-server) server listening on address: %s", addr)

		return nil
	}
}

func stopHook(rpcServer contracts.RPCServer) waiterContracts.CleanupFunc {
	return func(ctx context.Context) error {
		rpcServer.Logger().Info("(rpc-server) server shutting down...")

		if err := rpcServer.Shutdown(ctx); err != nil {
			rpcServer.Logger().Errorf("(rpc-server) failed to stop with error: %v", err)
			return nil
		}

		rpcServer.Logger().Info("(rpc-server) server shutdown done successfully")

		return nil
	}
}
