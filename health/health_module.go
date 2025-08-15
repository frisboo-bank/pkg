package health

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
)

var Module = module.NewModule(
	"health",

	provider.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.EnvConfig, error) {
			return config.LoadEnvConfig(loader, env)
		},
	),

	provider.Provide(
		func(
			loggerEnvCfg *loggerConfig.EnvConfig,
			envCfg *config.EnvConfig,
			httpServer httpServerContracts.HTTPServer,
			// services []contracts.HealthServiceCheck,
		) (contracts.HealthService, contracts.HealthEndpoint, error) {
			loggerOpts := loggerConfig.FromEnvConfig(loggerEnvCfg).
				With(loggerConfig.Prefix("health"))

			logger, err := logger.GetInstance(loggerEnvCfg.Type, loggerOpts)
			if err != nil {
				return nil, nil, err
			}

			opts := config.FromEnvConfig(envCfg)

			service := NewHealthService(logger, opts)
			// service.AddServices(services...)

			return service, NewHealthEndpoint(logger, httpServer, service, opts), nil
		}),

	invoker.Invoke(func(healthEndpoint contracts.HealthEndpoint) {
		healthEndpoint.RegisterEndpoints()
	}),
)
