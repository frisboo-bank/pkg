package httpserver

import (
	"context"

	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/options"
	"frisboo-bank/pkg/logger"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerOptions "frisboo-bank/pkg/logger/options"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

type HTTPServerDeps struct {
	dig.In
	Logger  loggerContracts.Logger `name:http_server_logger`
	Options *options.HTTPServerOptions
}

var Module = container.NewModule(
	"http_server",

	// load httpserver config
	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.HTTPServerOptions, error) {
			return options.ProvideHTTPServerOptions(loader, env)
		},
	),

	// create a custom httpserver logger
	container.Provide(func(options *loggerOptions.LogOptions) (loggerContracts.Logger, error) {
		logger, err := logger.GetInstance(options.Type)
		if err != nil {
			return nil, err
		}

		return logger.WithOptions(options).
			WithPrefix("http-server"), nil
	}, dig.Name("http_server_logger")),

	// create the httpserver
	container.Provide(func(deps HTTPServerDeps) (contracts.HTTPServer, error) {
		httpServer, err := GetInstance(deps.Options.Type, deps.Logger)
		if err != nil {
			return nil, err
		}

		return httpServer.WithOptions(deps.Options), nil
	}),

	container.Hook(startHook, stopHook),
)

func startHook(httpServer contracts.HTTPServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		httpServer.SetupDefaultMiddlewares()

		addr := httpServer.Address()

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

func stopHook(httpServer contracts.HTTPServer, logger loggerContracts.Logger) waiterContracts.CleanupFunc {
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
