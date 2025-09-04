package application

import (
	"context"

	"frisboo-bank/pkg/application/contracts"
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"
)

var _ contracts.Application = (*application)(nil)

type application struct {
	container   containerContracts.Container
	decorators  []decorator.Decorator
	environment environment.Environment
	hooks       []hook.Hooks
	invokers    []invoker.Invoker
	logger      loggerContracts.Logger
	modules     []module.Module
	providers   []provider.Provider
}

func NewApplication(
	modules []module.Module,
	providers []provider.Provider,
	decorators []decorator.Decorator,
	container containerContracts.Container,
	logger loggerContracts.Logger,
	environment environment.Environment,
) contracts.Application {
	validation.Assert(container != nil, "the container can't be nil", "application")
	validation.Assert(logger != nil, "the logger can't be nil", "application")

	return &application{
		container:   container,
		decorators:  decorators,
		environment: environment,
		logger:      logger,
		modules:     modules,
		providers:   providers,
	}
}

func (a *application) ResolveFunc(invoke invoker.Invoker) {
	a.invokers = append(a.invokers, invoke)
}

func (a *application) RegisterHook(hook hook.Hooks) {
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
		a.logger.Fatal("failed to stop because application not started")
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

	return a.container.RegisterModule(module.ModuleFunc("app",
		deps...,
	))
}

func (a *application) Logger() loggerContracts.Logger {
	return a.logger
}
