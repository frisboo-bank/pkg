package contracts

import (
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type ApplicationBuilder interface {
	ProvideModule(modules ...container.Module)
	Provide(providers ...container.Provider)
	Decorate(decorators ...container.Decorator)
	Build() Application

	GetModules() []container.Module
	GetProviders() []container.Provider
	GetDecorators() []container.Decorator
	Logger() loggerContracts.Logger
	Environment() environment.Environment
}
