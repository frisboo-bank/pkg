package contracts

import (
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type ApplicationBuilder interface {
	ProvideModule(modules ...dependencies.Module)
	Build() Application
	Modules() []dependencies.Module
	Providers() []dependencies.Provider
	Decorators() []dependencies.Decorator
	Container() containerContracts.Container
	Logger() loggerContracts.Logger
	Environment() environment.Environment
}
