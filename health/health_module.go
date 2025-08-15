package health

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
)

var Module = dependencies.NewModule(
	"health",

	dependencies.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.EnvConfig, error) {
			return config.LoadEnvConfig(loader, env)
		},
	),

	dependencies.Provide(
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

	dependencies.Invoke(func(healthEndpoint contracts.HealthEndpoint) {
		healthEndpoint.RegisterEndpoints()
	}),
)
