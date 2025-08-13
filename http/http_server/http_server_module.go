package httpserver

import (
	"context"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"

	configContracts "frisboo-bank/pkg/config/contracts"

	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = container.NewModule(
	"http_server",

	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.HTTPServerConfig, error) {
			return config.ProvideHTTPServerConfig(loader, env)
		},
	),

	container.Provide(func(loggerCfg *loggerConfig.LoggerConfig, options *config.HTTPServerConfig) (contracts.HTTPServer, error) {
		customLogger, err := logger.GetInstanceFromConfig(loggerCfg)
		if err != nil {
			return nil, err
		}
		customLogger = customLogger.WithPrefix("http-server")

		httpServer, err := GetInstanceFromConfig(options, customLogger)
		if err != nil {
			return nil, err
		}

		return httpServer, nil
	},
	),

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
