package goose

import (
	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
)

var _ contracts.MigratorAdapter = (*gooseMigratorAdapter)(nil)

type gooseMigratorAdapter struct {
	name   string
	cfg    *config.Config
	client databaseclientContracts.DatabaseClient
	logger loggerContracts.Logger
}

func New(name string, cfg *config.Config, dbClient databaseclientContracts.DatabaseClient, logger loggerContracts.Logger) *gooseMigratorAdapter {
	return &gooseMigratorAdapter{
		name:   name,
		cfg:    cfg,
		client: dbClient,
		logger: logger,
	}
}

func (g *gooseMigratorAdapter) Run(commandType migrationcommandtype.MigrationCommandType, version string) error {
	panic("unimplemented")
}

func (g *gooseMigratorAdapter) Config() *config.Config {
	return g.cfg
}

func (g *gooseMigratorAdapter) Logger() loggerContracts.Logger {
	return g.logger
}

func (g *gooseMigratorAdapter) Name() string {
	return g.name
}

func (g *gooseMigratorAdapter) Type() migratortype.MigratorType {
	return migratortype.MigratorTypes.GOOSE
}
