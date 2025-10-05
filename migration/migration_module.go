package migration

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	"frisboo-bank/pkg/logger"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	loggerConfig "frisboo-bank/pkg/logger/config"
)

const MigrationsGroup = "migrations"

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	validation.AssertNotNil("appBuilder", appBuilder)

	configLoader := appBuilder.ConfigLoader()
	env := appBuilder.Environment()
	logger := appBuilder.Logger()

	m := module.ModuleFunc("migration")

	// Load and register the config registry
	cfgRegistry, err := config.LoadRegistry(configLoader, env)
	if err != nil {
		logger.Fatalf("failed to register migration module with error: %v", err)
	}
	m.AddProvider(provider.ProvideFunc(func() config.Registry { return cfgRegistry }))

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			logger.Fatalf("failed to register migration:{%s} module with error:{%v}", name, err)
		}
		m.AddModule(serverModuleFunc(name, logger, &cfg))
	}

	return m
}

func serverModuleFunc(name string, log loggerContracts.Logger, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("log", log)
	validation.AssertNotNil("cfg", cfg)

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
		provider.Name("migration:"+name),
		provider.Group(MigrationsGroup),
		provider.NamedDep("dbClientRef", "database-client:"+name),
	))

	return m
}
