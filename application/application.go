package application

import (
	"context"
	"fmt"
	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"os"
)

type application struct {
	modules     []container.Module
	providers   []container.Provider
	decorators  []container.Decorator
	invokers    []container.Invoker
	hooks       []container.HookStarter
	container   contracts.ApplicationContainer
	logger      loggerContracts.Logger
	environment environment.Environment
}

var _ contracts.Application = (*application)(nil)

func NewApplication(
	modules []container.Module,
	providers []container.Provider,
	decorators []container.Decorator,
	logger loggerContracts.Logger,
	environment environment.Environment,
) contracts.Application {
	return &application{
		modules:     modules,
		providers:   providers,
		decorators:  decorators,
		logger:      logger,
		environment: environment,
	}
}

func (a *application) ResolveFunc(invoke container.Invoker) {
	a.invokers = append(a.invokers, invoke)
}

func (a *application) RegisterHook(hook container.HookStarter) {
	a.hooks = append(a.hooks, hook)
}

func (a *application) Start(ctx context.Context) error {
	container := NewApplicationContainer(a)
	a.container = container

	return container.Start(ctx)
}

func (a *application) Stop(ctx context.Context) error {
	if a.container == nil {
		fmt.Println("application: Failed to stop because application not started.")
		os.Exit(1)
	}
	return a.container.Stop(ctx)
}

func (a *application) Logger() loggerContracts.Logger {
	return a.logger
}

func (a *application) Environment() environment.Environment {
	return a.environment
}
