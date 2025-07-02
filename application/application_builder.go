package application

import (
	"fmt"
	"os"

	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/application/infrastructure"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerFactory "frisboo-bank/pkg/logger/factory"
	loggerOptions "frisboo-bank/pkg/logger/options"
)

type applicationBuilder struct {
	providers   []container.Provider
	decorators  []container.Decorator
	modules     []container.Module
	logger      loggerContracts.Logger
	environment environment.Environment
}

var _ contracts.ApplicationBuilder = (*applicationBuilder)(nil)

func NewApplicationBuilder(environments ...environment.Environment) contracts.ApplicationBuilder {
	env := environment.GetEnvFromConfig(environments...)

	logOptions, err := loggerOptions.ProvideLogOptions(env)
	if err != nil {
		fmt.Printf("application-builder: failed to load Logger options with error: %v\n", err)
		os.Exit(1)
	}

	logger, err := loggerFactory.GetInstance(logOptions)
	if err != nil {
		fmt.Printf("application-builder: failed to create Logger with error: %v\n", err)
		os.Exit(1)
	}

	return &applicationBuilder{
		logger:      logger,
		environment: env,
		modules: []container.Module{
			infrastructure.Module,
		},
		providers: []container.Provider{
			container.Provide(func() environment.Environment { return env }),
			container.Provide(func() loggerContracts.Logger { return logger }),
		},
	}
}

func (b *applicationBuilder) ProvideModule(modules ...container.Module) {
	b.modules = append(b.modules, modules...)
}

func (b *applicationBuilder) Provide(providers ...container.Provider) {
	b.providers = append(b.providers, providers...)
}

func (b *applicationBuilder) Decorate(decorators ...container.Decorator) {
	b.decorators = append(b.decorators, decorators...)
}

func (b *applicationBuilder) Build() contracts.Application {
	return NewApplication(
		b.modules,
		b.providers,
		b.decorators,
		b.logger,
		b.environment,
	)
}

func (b *applicationBuilder) GetModules() []container.Module {
	return b.modules
}

func (b *applicationBuilder) GetProviders() []container.Provider {
	return b.providers
}

func (b *applicationBuilder) GetDecorators() []container.Decorator {
	return b.decorators
}

func (b *applicationBuilder) Logger() loggerContracts.Logger {
	return b.logger
}

func (b *applicationBuilder) Environment() environment.Environment {
	return b.environment
}
