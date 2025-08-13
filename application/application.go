package application

import (
	"context"
	"fmt"
	"os"

	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"

	containerContracts "frisboo-bank/pkg/container/contracts"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type application struct {
	container   containerContracts.Container
	decorators  []containerContracts.Decorator
	environment environment.Environment
	hooks       []containerContracts.HookStarter
	invokers    []containerContracts.Invoker
	logger      loggerContracts.Logger
	modules     []containerContracts.Module
	providers   []containerContracts.Provider
}

var _ contracts.Application = (*application)(nil)

func NewApplication(
	modules []containerContracts.Module,
	providers []containerContracts.Provider,
	decorators []containerContracts.Decorator,
	container containerContracts.Container,
	logger loggerContracts.Logger,
	environment environment.Environment,
) contracts.Application {
	return &application{
		modules:     modules,
		providers:   providers,
		decorators:  decorators,
		environment: environment,
		container:   container,
		logger:      logger,
	}
}

func (a *application) ResolveFunc(invoke containerContracts.Invoker) {
	a.invokers = append(a.invokers, invoke)
}

func (a *application) RegisterHook(hook containerContracts.HookStarter) {
	a.hooks = append(a.hooks, hook)
}

func (a *application) Start(ctx context.Context) error {
	if err := a.registerDependencies(); err != nil {
		return err
	}

	return a.container.Start(ctx)
}

func (a *application) Stop(ctx context.Context) error {
	if a.container == nil {
		fmt.Println("application: Failed to stop because application not started.")
		os.Exit(1)
	}
	return a.container.Stop(ctx)
}

func (a *application) registerDependencies() error {
	dependenciesLen := len(a.modules) + len(a.providers) + len(a.decorators) + len(a.invokers) + len(a.hooks)
	dependencies := make([]containerContracts.Dependency, 0, dependenciesLen)

	for _, dep := range a.modules {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range a.providers {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range a.decorators {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range a.invokers {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range a.hooks {
		dependencies = append(dependencies, dep)
	}

	return a.container.RegisterModule(container.NewModule("app",
		dependencies...,
	))
}
