package rpcserver

import (
	"context"
	"errors"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	"net"

	configContracts "frisboo-bank/pkg/config/contracts"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerOptions "frisboo-bank/pkg/logger/options"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type RPCServerDeps struct {
	dig.In
	Logger  loggerContracts.Logger `name:rpcServerLogger`
	Options *options.RPCServerOptions
}

var Module = container.NewModule(
	"rpc-server",

	// load rpc config
	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.RPCServerOptions, error) {
			return options.ProvideRPCServerOptions(loader, env)
		},
	),

	// create a custom rpcserver logger
	container.Provide(func(options *loggerOptions.LogOptions) (loggerContracts.Logger, error) {
		logger, err := logger.GetInstance(options.Type)
		if err != nil {
			return nil, err
		}

		logger.WithOptions(options).
			WithPrefix("rpc-server")

		return logger, nil
	}, dig.Name("rpcServerLogger")),

	// create the rpcserver
	container.Provide(
		func(deps RPCServerDeps) (contracts.RPCServer, error) {
			rpcServer, err := GetInstance(deps.Logger)
			if err != nil {
				return nil, err
			}

			return rpcServer.WithOptions(deps.Options), nil
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
