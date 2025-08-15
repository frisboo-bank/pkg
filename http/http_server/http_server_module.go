package httpserver

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"

	configContracts "frisboo-bank/pkg/config/contracts"

	loggerConfig "frisboo-bank/pkg/logger/config"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = module.NewModule(
	"http_server",

	provider.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.EnvConfig, error) {
			return config.LoadEnvConfig(loader, env)
		},
	),

	provider.Provide(
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

	hook.Hook(func(httpServer contracts.HTTPServer) waiterContracts.WaitFunc {
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
