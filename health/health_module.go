package health

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
)

var Module = container.NewModule(
	"health",
	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*config.HealthConfig, error) {
			return config.ProvideHealthConfig(loader, env)
		},
	),
	container.Provide(func(config *config.HealthConfig) contracts.HealthService {
		return NewHealthService(nil)
	}),
	container.Provide(
		func(
			config *config.HealthConfig,
			healthService contracts.HealthService,
			httpServer httpServerContracts.HTTPServer,
		) contracts.HealthEndpoint {
			return NewHealthEndpoint(
				healthService,
				httpServer,
			)
		},
	),
	container.Invoke(func(healthEndpoint contracts.HealthEndpoint) {
		healthEndpoint.RegisterEndpoints()
	}),
)
