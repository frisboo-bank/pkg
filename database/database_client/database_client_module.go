package databaseclient

import (
	"context"
	"fmt"

	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/database/database_client/contracts"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

const (
	DatabaseClientsGroup   = "database-clients"
	DatabaseClientProvider = "database-client:%s"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	validation.AssertNotNil("appBuilder", appBuilder)

	configLoader := appBuilder.ConfigLoader()
	env := appBuilder.Environment()
	logger := appBuilder.Logger()

	m := module.ModuleFunc("database-clients")

	// Load and register the config registry
	cfgRegistry, err := config.LoadRegistry(configLoader, env)
	if err != nil {
		logger.Panicw("failed to register database_client module", loggerContracts.Fields{"err": err, "cause": syserrors.Cause(err)})
	}
	m.AddProvider(provider.ProvideFunc(func() config.Registry { return cfgRegistry }))

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			logger.Panicw("failed to register database_client module", loggerContracts.Fields{"err": err, "cause": syserrors.Cause(err)})
		}
		if !cfg.Enabled {
			logger.Debugf("database-client:{%s} is disabled and will not be loaded", name)
			continue
		}
		m.AddModule(serverModuleFunc(name, logger, &cfg))
	}

	return m
}

func serverModuleFunc(name string, log loggerContracts.Logger, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("log", log)
	validation.AssertNotNil("cfg", cfg)

	log.Debugf("Try to register database-client:{%s} module", name)

	m := module.ModuleFunc("database-client:" + name)

	type providerProps struct {
		LoggerCfgRegistry loggerConfig.Registry
		AppLogger         loggerContracts.Logger
	}

	m.AddProvider(provider.ProvideFunc(func(props providerProps) (contracts.DatabaseClient, error) {
		loggerCfgRegistry := props.LoggerCfgRegistry
		appLogger := props.AppLogger

		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "database-client:{%s} logger", name)
		}
		return GetInstance(name, cfg, log)
	},
		provider.Name(fmt.Sprintf(DatabaseClientProvider, name)),
		provider.Group(DatabaseClientsGroup),
	))

	type hookParams struct {
		DatabaseClient contracts.DatabaseClient `name:"dbClientRef"`
	}

	m.AddHook(hook.HooksFunc(fmt.Sprintf("database-client-%s-hook", name),
		func(p hookParams) waiterContracts.WaitFunc {
			return func(ctx context.Context) error {
				clt := p.DatabaseClient

				log.Debugf("Try to ping database:{%s}", clt.Name())

				if err := clt.Ping(); err != nil {
					clt.Logger().Errorf("failed to ping the database:{%s} with error:{%v}", clt.Name(), err)
					return err
				}
				clt.Logger().Infof("database:{%s} pinged successfully", clt.Name())

				return nil
			}
		},
		func(p hookParams) waiterContracts.CleanupFunc {
			return func(ctx context.Context) error {
				clt := p.DatabaseClient

				if err := clt.Disconnect(); err != nil {
					clt.Logger().Errorf("disconnecting from database:{%s} failed with error:{%w}", clt.Name(), err)
				} else {
					clt.Logger().Infof("disconnected from database:{%s} successfully", clt.Name())
				}

				return nil
			}
		},
		hook.NamedDep("dbClientRef", fmt.Sprintf(DatabaseClientProvider, name)),
	))

	return m
}
