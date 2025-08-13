package httpserver

import (
	"context"

	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = container.NewModule(
	"http_server",

	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.EnvConfig, error) {
			return config.LoadEnvConfig(loader, env)
		},
	),

	container.Provide(
		func(loggerEnvCfg *loggerConfig.EnvConfig, envCfg *config.EnvConfig) (contracts.HTTPServer, error) {
			loggerOpts := loggerConfig.FromEnvConfig(loggerEnvCfg).
				With(loggerConfig.Prefix("http-server"))

			logger, err := logger.GetInstance(loggerEnvCfg.Type, loggerOpts)
			if err != nil {
				return nil, err
			}

			opts := config.FromEnvConfig(envCfg)

			httpServer, err := GetInstance(envCfg.Type, logger, opts)
			if err != nil {
				return nil, err
			}

			return httpServer, nil
		},
	),

	container.Hook(func(httpServer contracts.HTTPServer) waiterContracts.WaitFunc {
		return func(ctx context.Context) error {
			httpServer.SetupDefaultMiddlewares()

			var err error
			go func() {
				err = httpServer.Start(ctx)
			}()

			return err
		}
	}, func(httpServer contracts.HTTPServer) waiterContracts.CleanupFunc {
		return func(ctx context.Context) error {
			return httpServer.Shutdown(ctx)
		}
	}),
)
