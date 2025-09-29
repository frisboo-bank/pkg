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

const DatabaseClientsGroup = "database-clients"

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	m := module.ModuleFunc("database-clients")

	// Load and register the config registry
	cfgRegistry, err := config.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
	if err != nil {
		appBuilder.Logger().Fatalf("failed to register database-clients module with error: %v", err)
	}
	m.AddProvider(provider.ProvideFunc(func() config.Registry { return cfgRegistry }))

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			appBuilder.Logger().Fatal("failed to register database-client:{%s} module with error:{%v}", name, err)
		}
		if !cfg.Enabled {
			appBuilder.Logger().Debugf("database-client:{%s} is disabled and will not be loaded", name)
			continue
		}
		m.AddModule(serverModuleFunc(name, appBuilder.Logger(), &cfg))
	}

	return m
}

func serverModuleFunc(name string, log loggerContracts.Logger, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("log", log)
	validation.AssertNotNil("cfg", cfg)

	log.Debugf("Try to register database-client:{%s} module", name)

	m := module.ModuleFunc("database-client:" + name)

	// Instance registration name
	providerName := "database-client:" + name

	m.AddProvider(provider.ProvideFunc(func(loggerCfgRegistry loggerConfig.Registry, appLogger loggerContracts.Logger) (contracts.DatabaseClient, error) {
		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "http-server:{%s} logger", name)
		}
		return GetInstance(name, cfg, log)
	},
		provider.Name(providerName),
		provider.Group(DatabaseClientsGroup),
	))

	type hookParams struct {
		DatabaseClient contracts.DatabaseClient `name:"databaseClientRef"`
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
					clt.Logger().Errorf("disconnecting from database:{%s} failed with error:{%v}", clt.Name(), err)
				} else {
					clt.Logger().Infof("disconnected from database:{%s} successfully", clt.Name())
				}

				return nil
			}
		},
		hook.NamedDep("databaseClientRef", providerName),
	))

	return m
}
