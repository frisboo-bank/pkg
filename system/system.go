package system

import (
	"fmt"
	"frisboo-bank/pkg/di"
	diConfig "frisboo-bank/pkg/di/config"
	diFactory "frisboo-bank/pkg/di/factory"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerFactory "frisboo-bank/pkg/logger/factory"
	"frisboo-bank/pkg/system/contracts"
	"frisboo-bank/pkg/waiter"
)

type system struct {
	di          di.Container
	environment environment.Environment
	logger      logger.Logger
	waiter      waiter.Waiter
}

func NewSystem(env environment.Environment) contracts.System {
	system := &system{
		environment: env,
	}

	loggerCfg, err := loggerConfig.ProvideLogConfig(env)
	if err != nil {
		panic(fmt.Errorf("system: failed to load the logger config with error: %w", err))
	}

	system.logger = loggerFactory.NewInstance(loggerCfg)

	system.di = diFactory.NewInstance(&diConfig.DiOptions{
		Type:   di.DiTypeSimpleDi,
		Logger: system.logger,
	})

	return system
}

func (a *system) Run() error {
	return nil
}

func (a *system) Di() di.Container {
	return a.di
}

func (a *system) Logger() logger.Logger {
	return a.logger
}

func (a *system) Environment() environment.Environment {
	return a.environment
}
