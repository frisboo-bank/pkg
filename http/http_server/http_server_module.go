package httpserver

import (
	"context"

	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/factory"
	"frisboo-bank/pkg/http/http_server/options"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = container.NewModule(
	"http_server",
	container.Provide(func(env environment.Environment) (*options.HttpServerOptions, error) {
		return options.ProvideHttpServerOptions(env)
	}),
	container.Provide(
		func(config *options.HttpServerOptions, logger loggerContracts.Logger) (contracts.HttpServer, error) {
			logger = logger.WithName("http-server")

			return factory.GetInstance(config,
				options.WithLogger(logger),
			)
		},
	),
	container.Hook(startHook, stopHook),
)

func startHook(httpServer contracts.HttpServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		httpServer.SetupDefaultMiddlewares()

		addr := httpServer.Config().Address()

		httpServer.Logger().Info("(http-server) starting server...")

		go func() {
			if err := httpServer.Start(); err != nil {
				httpServer.Logger().Fatalf("(http-server) failed to start with error: %v", err)
			}
		}()

		httpServer.Logger().Infof("(http-server) server listening on address: %s", addr)

		return nil
	}
}

func stopHook(httpServer contracts.HttpServer, logger loggerContracts.Logger) waiterContracts.CleanupFunc {
	return func(ctx context.Context) error {
		httpServer.Logger().Info("(http-server) server shutting down...")

		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Logger().Errorf("(http-server) failed to stop with error: %v", err)
			return nil
		}

		httpServer.Logger().Info("(http-server) server shutdown done successfully")

		return nil
	}
}
