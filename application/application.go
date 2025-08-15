package application

import (
	"context"
	"fmt"
	"os"

	"frisboo-bank/pkg/application/contracts"
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type application struct {
	container   containerContracts.Container
	decorators  []dependencies.Decorator
	environment environment.Environment
	hooks       []dependencies.Hooks
	invokers    []dependencies.Invoker
	logger      loggerContracts.Logger
	modules     []dependencies.Module
	providers   []dependencies.Provider
}

var _ contracts.Application = (*application)(nil)

func NewApplication(
	modules []dependencies.Module,
	providers []dependencies.Provider,
	decorators []dependencies.Decorator,
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

func (a *application) ResolveFunc(invoke dependencies.Invoker) {
	a.invokers = append(a.invokers, invoke)
}

func (a *application) RegisterHook(hook dependencies.Hooks) {
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
	deps := make([]dependencies.Dependency, 0, dependenciesLen)

	for _, dep := range a.modules {
		deps = append(deps, dep)
	}

	for _, dep := range a.providers {
		deps = append(deps, dep)
	}

	for _, dep := range a.decorators {
		deps = append(deps, dep)
	}

	for _, dep := range a.invokers {
		deps = append(deps, dep)
	}

	for _, dep := range a.hooks {
		deps = append(deps, dep)
	}

	return a.container.RegisterModule(dependencies.NewModule("app",
		deps...,
	))
}
