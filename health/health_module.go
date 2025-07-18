package health

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/health/options"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
)

var Module = container.NewModule(
	"health",
	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.HealthOptions, error) {
			return options.ProvideHealthOptions(loader, env)
		},
	),
	container.Provide(func(config *options.HealthOptions) contracts.HealthService {
		return NewHealthService(config.Services)
	}),
	container.Provide(
		func(
			config *options.HealthOptions,
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
