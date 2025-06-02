package application

import (
	"errors"
	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/di/simple"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	loggerDi "frisboo-bank/pkg/logger/di"
	"frisboo-bank/pkg/waiter"

	customErrors "frisboo-bank/pkg/custom_errors"
)

var ApplicationFailedToStartError = customErrors.WrapAsFatal(errors.New("application: failed to start"))

type application struct {
	di          di.Container
	environment environment.Environment
	logger      logger.Logger
	waiter      waiter.Waiter
}

func NewApplication(env environment.Environment) (contracts.Application, error) {
	app := &application{
		di:          simple.NewSimpleDi(),
		environment: env,
		waiter:      waiter.NewWaiter(waiter.CatchSignals()),
	}

	app.di.AddSingleton("logger", loggerDi.LoggerFactory)

	return app, nil
}

func (a *application) Run() error {
	return a.Waiter().Wait()
}

func (a *application) Di() di.Container {
	return a.di
}

func (a *application) Waiter() waiter.Waiter {
	return a.waiter
}

func (a *application) Logger() logger.Logger {
	return a.logger
}

func (a *application) Environment() environment.Environment {
	return a.environment
}
