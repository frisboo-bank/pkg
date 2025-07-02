package application

import (
	"context"
	"fmt"

	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container"
	containerContracts "frisboo-bank/pkg/container/contracts"
	digContainer "frisboo-bank/pkg/container/dig"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type applicationContainer struct {
	modules   []container.Module
	container containerContracts.Container
	logger    loggerContracts.Logger
}

var _ contracts.ApplicationContainer = (*applicationContainer)(nil)

func NewApplicationContainer(app *application) contracts.ApplicationContainer {
	dependenciesLen := len(app.modules) + len(app.providers) + len(app.decorators) + len(app.invokers) + len(app.hooks)
	dependencies := make([]container.Dependency, 0, dependenciesLen)

	for _, dep := range app.modules {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range app.providers {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range app.decorators {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range app.invokers {
		dependencies = append(dependencies, dep)
	}

	for _, dep := range app.hooks {
		dependencies = append(dependencies, dep)
	}

	appModule := container.NewModule("app",
		dependencies...,
	)

	return &applicationContainer{
		modules: []container.Module{
			appModule,
		},
		logger: app.logger,
	}
}

func (a *applicationContainer) Start(ctx context.Context) error {
	container := digContainer.NewDigContainer(a.modules, a.logger)
	a.container = container

	return a.container.Start(ctx)
}

func (a *applicationContainer) Stop(ctx context.Context) error {
	if a.container == nil {
		return fmt.Errorf("application container: looks like there is no container currently running")
	}

	return a.container.Stop(ctx)
}

func (a *applicationContainer) Logger() loggerContracts.Logger {
	return a.logger
}
