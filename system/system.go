package system

import (
	"context"
	"fmt"
	"frisboo-bank/pkg/di"
	diConfig "frisboo-bank/pkg/di/config"
	diFactory "frisboo-bank/pkg/di/factory"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerFactory "frisboo-bank/pkg/logger/factory"
	"frisboo-bank/pkg/system/contracts"

	"github.com/jackc/pgx/v5"
)

type system struct {
	di          di.Container
	environment environment.Environment
	logger      logger.Logger
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

	system.di.AddSingleton("database", func(r di.Resolver) (any, error) {
		return pgx.Connect(context.Background(), "user=pqgotest dbname=pqgotest sslmode=verify-full")
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
