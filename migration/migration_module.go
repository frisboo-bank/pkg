package migration

import (
	"fmt"

	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

const (
	MigrationsGroup    = "migrations"
	MigrationsProvider = "migration:%s"
)

type ModuleProps struct {
	AppBuilder  applicationContracts.ApplicationBuilder
	CfgRegistry config.Registry
}

func ModuleFunc(props ModuleProps) module.Module {
	validation.AssertNotNil("props", props)
	validation.AssertNotNil("props.AppBuilder", props.AppBuilder)

	appBuilder := props.AppBuilder
	logger := appBuilder.Logger()
	cfgRegistry := props.CfgRegistry

	if cfgRegistry == nil {
		var err error
		cfgRegistry, err = config.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
		if err != nil {
			logger.Panicf("failed to register migration module with error: %v", err)
		}
	}

	m := module.ModuleFunc(
		"migration",
		provider.ProvideFunc(func() config.Registry { return cfgRegistry }),
	)

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			logger.Panicf("failed to register migration:{%s} module with error:{%v}", name, err)
		}
		m.AddModule(serverModuleFunc(name, &cfg, logger))
	}

	return m
}

func serverModuleFunc(name string, cfg *config.Config, log loggerContracts.Logger) module.Module {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("log", log)

	log.Debugf("Try to register migration:{%s} module", name)

	m := module.ModuleFunc("migration:" + name)

	type providerParams struct {
		DBClient          databaseclientContracts.DatabaseClient `name:"dbClientRef"`
		LoggerCfgRegistry loggerConfig.Registry
		AppLogger         loggerContracts.Logger
	}

	m.AddProvider(provider.ProvideFunc(func(props providerParams) (contracts.Migrator, error) {
		dbClient := props.DBClient
		loggerCfgRegistry := props.LoggerCfgRegistry
		appLogger := props.AppLogger

		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "migration:{%s} logger", name)
		}
		return GetInstance(name, cfg, dbClient, log)
	},
		provider.Name(fmt.Sprintf(MigrationsProvider, name)),
		provider.Group(MigrationsGroup),
		provider.NamedDep("dbClientRef", "database-client:"+name),
	))

	return m
}
