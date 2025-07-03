package rpcserver

import (
	"context"
	"errors"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/factory"
	"frisboo-bank/pkg/rpc/rpc_server/options"
	"net"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
	"google.golang.org/grpc"
)

type RPCServerDeps struct {
	dig.In
	Config *options.RPCServerOptions
	Logger loggerContracts.Logger `name:"rpcServerLogger"`
}

var Module = container.NewModule(
	"rpc-server",

	container.Provide(func(env environment.Environment) (*options.RPCServerOptions, error) {
		return options.ProvideRPCServerOptions(env)
	}),

	container.Provide(func(logger loggerContracts.Logger) loggerContracts.Logger {
		return logger.
			WithPrefix("(rpc-Server)").
			WithName("rpc-server")
	}, dig.Name("rpcServerLogger")),

	container.Provide(
		func(dependencies RPCServerDeps) (contracts.RPCServer, error) {
			return factory.GetInstance(dependencies.Config,
				options.UseLogger(dependencies.Logger),
			)
		},
	),

	container.Hook(startHook, stopHook),
)

func startHook(rpcServer contracts.RPCServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		addr := rpcServer.Config().Address()

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
