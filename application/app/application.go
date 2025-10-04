package app

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
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

var _ contracts.Application = (*application)(nil)

type application struct {
	modules     []module.Module
	providers   []provider.Provider
	hooks       []hook.Hooks
	decorators  []decorator.Decorator
	invokers    []invoker.Invoker
	container   containerContracts.Container
	logger      loggerContracts.Logger
	environment environment.Environment
}

func NewApplication(
	container containerContracts.Container,
	logger loggerContracts.Logger,
	environment environment.Environment,
	modules []module.Module,
	providers []provider.Provider,
	decorators []decorator.Decorator,
) contracts.Application {
	validation.AssertNotNil("container", container)
	validation.AssertNotNil("logger", logger)
	validation.AssertNotNil("environment", environment)

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

func (a *application) Run() error {
	if err := a.registerDependencies(); err != nil {
		return err
	}
	return nil
}

func (a *application) Start(ctx context.Context) error {
	if err := a.registerDependencies(); err != nil {
		return err
	}
	return a.container.Start(ctx)
}

func (a *application) Stop(ctx context.Context) error {
	if a.container == nil {
		return syserrors.New("failed to stop because application not started")
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

func (a *application) Environment() environment.Environment {
	return a.environment
}
