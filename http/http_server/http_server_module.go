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

	"go.uber.org/dig"
)

type HttpServerDeps struct {
	dig.In
	Config *options.HttpServerOptions
	Logger loggerContracts.Logger `name:"httpServerLogger"`
}

var Module = container.NewModule(
	"http_server",
	container.Provide(func(env environment.Environment) (*options.HttpServerOptions, error) {
		return options.ProvideHttpServerOptions(env)
	}),

	container.Provide(func(logger loggerContracts.Logger) loggerContracts.Logger {
		return logger.
			WithPrefix("(http-Server)").
			WithName("http-server")
	}, dig.Name("httpServerLogger")),

	container.Provide(
		func(dependencies HttpServerDeps) (contracts.HttpServer, error) {
			return factory.GetInstance(dependencies.Config,
				options.UseLogger(dependencies.Logger),
			)
		},
	),
	container.Hook(startHook, stopHook),
)

func startHook(httpServer contracts.HttpServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		httpServer.SetupDefaultMiddlewares()

		addr := httpServer.Config().Address()

		httpServer.Logger().Info("starting server...")

		go func() {
			if err := httpServer.Start(); err != nil {
				httpServer.Logger().Fatalf("failed to start with error: %v", err)
			}
		}()

		httpServer.Logger().Infof("server listening on address: %s", addr)

		return nil
	}
}

func stopHook(httpServer contracts.HttpServer, logger loggerContracts.Logger) waiterContracts.CleanupFunc {
	return func(ctx context.Context) error {
		httpServer.Logger().Info("server shutting down...")

		if err := httpServer.Shutdown(ctx); err != nil {
			httpServer.Logger().Errorf("failed to stop with error: %v", err)
			return nil
		}

		httpServer.Logger().Info("server shutdown done successfully")

		return nil
	}
}
