package health

import (
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/health/endpoints"
	"frisboo-bank/pkg/health/options"
	"frisboo-bank/pkg/health/services"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
)

var Module = container.NewModule("health",
	container.Provide(func(env environment.Environment) (*options.HealthOptions, error) {
		return options.ProvideHealthOptions(env)
	}),
	container.Provide(func(config *options.HealthOptions) contracts.HealthService {
		return services.NewHealthCheckService(config)
	}),
	container.Provide(func(config *options.HealthOptions, healthService contracts.HealthService, httpServer httpServerContracts.HttpServer) contracts.HealthEndpoint {
		return endpoints.NewHealthCheckEndpoint(config, healthService, httpServer)
	}),
	container.Invoke(func(healthEndpoint contracts.HealthEndpoint) {
		healthEndpoint.RegisterEndpoints()
	}),
)
