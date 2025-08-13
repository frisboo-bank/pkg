package health

import (
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/logger"

	configContracts "frisboo-bank/pkg/config/contracts"

	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"

	loggerConfig "frisboo-bank/pkg/logger/config"
)

var Module = container.NewModule(
	"health",

	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.EnvConfig, error) {
			return config.LoadEnvConfig(loader, env)
		},
	),

	container.Provide(
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

	container.Invoke(func(healthEndpoint contracts.HealthEndpoint) {
		healthEndpoint.RegisterEndpoints()
	}),
)
