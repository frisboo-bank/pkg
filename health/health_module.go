package health

import (
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/logger"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"github.com/davecgh/go-spew/spew"
)

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"health",

		provider.ProvideFunc(config.Load),
		provider.ProvideFunc(func(cfg *config.Config) (loggerContracts.Logger, error) {
			return logger.GetInstance(&cfg.Logger)
		}, provider.Name("health_logger")),
		provider.ProvideFunc(NewHealthService),
		provider.ProvideFunc(NewHealthEndpoint),

		invoker.InvokerFunc(func(healthEndpoint contracts.HealthEndpoint) {
			spew.Dump(healthEndpoint)
			healthEndpoint.RegisterEndpoints()
		}),
	)
}
