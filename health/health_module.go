package health

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/validation"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	validation.AssertNotNil("appBuilder", appBuilder)

	configLoader := appBuilder.ConfigLoader()
	env := appBuilder.Environment()
	logger := appBuilder.Logger()

	_ = configLoader
	_ = env
	_ = logger

	m := module.ModuleFunc("health")

	return m

	// return module.ModuleFunc(
	// 	"health",
	//
	// 	provider.ProvideFunc(config.Load),
	// 	provider.ProvideFunc(func(cfg *config.Config) (loggerContracts.Logger, error) {
	// 		return logger.GetInstance(cfg.Logger)
	// 	}, provider.Name("health_logger")),
	// 	provider.ProvideFunc(NewHealthService),
	// 	provider.ProvideFunc(NewHealthEndpoint),
	//
	// 	invoker.InvokerFunc(func(healthEndpoint contracts.HealthEndpoint) {
	// 		healthEndpoint.RegisterEndpoints()
	// 	}),
	// )
}
