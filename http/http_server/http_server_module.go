package httpserver

import (
	"context"

	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/options"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = container.NewModule(
	"http_server",
	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.HTTPServerOptions, error) {
			return options.ProvideHTTPServerOptions(loader, env)
		},
	),

	container.Provide(
		func(logger loggerContracts.Logger, config *options.HTTPServerOptions) (contracts.HTTPServer, error) {
			httpServerLogger := logger.Clone()
			httpServerLogger.WithPrefix("http-server")

			httpServer, err := GetInstance(config.Type, httpServerLogger)
			if err != nil {
				return nil, err
			}

			return httpServer.
				WithBasePath(config.BasePath).
				WithBodyLimit(config.BodyLimit).
				HasDevelopment(config.Development).
				WithHost(config.Host).
				WithIdleTimeout(config.IdleTimeout).
				WithIgnoreLogUrls(config.IgnoreLogUrls).
				WithMaxHeaderBytes(config.MaxHeaderBytes).
				WithPort(config.Port).
				WithReadHeaderTimeout(config.ReadHeaderTimeout).
				WithReadTimeout(config.ReadTimeout).
				WithServerShutdownTimeout(config.ServerShutdownTimeout).
				WithWriteTimeout(config.WriteTimeout), nil
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
