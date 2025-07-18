package rpcserver

import (
	"context"
	"errors"
	"net"

	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"google.golang.org/grpc"
)

var Module = container.NewModule(
	"rpc-server",

	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.RPCServerOptions, error) {
			return options.ProvideRPCServerOptions(loader, env)
		},
	),

	container.Provide(
		func(logger loggerContracts.Logger, config *options.RPCServerOptions) (contracts.RPCServer, error) {
			rpcServerLogger := logger.Clone()
			rpcServerLogger.WithPrefix("rpc-server")

			rpcServer, err := GetInstance(rpcServerLogger)
			if err != nil {
				return nil, err
			}

			return rpcServer.
				WithHost(config.Host).
				WithPort(config.Port).
				WithServerShutdownTimeout(config.ServerShutdownTimeout), nil
		},
	),

	container.Hook(startHook, stopHook),
)

func startHook(rpcServer contracts.RPCServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		addr := rpcServer.Address()

		rpcServer.Logger().Info("starting server...")

		go func() {
			listener, err := net.Listen("tcp", addr)
			if err != nil {
				rpcServer.Logger().Fatalf("failed to listen on address with error: %v", err)
				return
			}

			if err := rpcServer.Start(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
				rpcServer.Logger().Fatalf("failed to start with error: %v", err)
			}
		}()

		rpcServer.Logger().Infof("server listening on address: %s", addr)

		return nil
	}
}

func stopHook(rpcServer contracts.RPCServer) waiterContracts.CleanupFunc {
	return func(ctx context.Context) error {
		rpcServer.Logger().Info("server shutting down...")

		if err := rpcServer.Shutdown(ctx); err != nil {
			rpcServer.Logger().Errorf("failed to stop with error: %v", err)
			return nil
		}

		rpcServer.Logger().Info("server shutdown done successfully")

		return nil
	}
}
