package contracts

import (
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type ApplicationBuilder interface {
	ProvideModule(modules ...containerContracts.Module)
	Build() Application
	Modules() []containerContracts.Module
	Providers() []containerContracts.Provider
	Decorators() []containerContracts.Decorator
	Container() containerContracts.Container
	Logger() loggerContracts.Logger
	Environment() environment.Environment
}
