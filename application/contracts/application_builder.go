package contracts

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type ApplicationBuilder interface {
	ProvideModule(modules ...module.Module)
	Build() Application
	Modules() []module.Module
	Providers() []provider.Provider
	Decorators() []decorator.Decorator
	Container() containerContracts.Container
	Logger() loggerContracts.Logger
	ConfigLoader() configloaderContracts.ConfigLoader
	Environment() environment.Environment
}
