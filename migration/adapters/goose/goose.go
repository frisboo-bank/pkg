package goose

import (
	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/validation"
)

var _ contracts.MigratorAdapter = (*gooseMigratorAdapter)(nil)

type gooseMigratorAdapter struct {
	name     string
	cfg      *config.Config
	dbClient databaseclientContracts.DatabaseClient
	logger   loggerContracts.Logger
}

func New(
	name string,
	cfg *config.Config,
	dbClient databaseclientContracts.DatabaseClient,
	logger loggerContracts.Logger,
) *gooseMigratorAdapter {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("dbClient", dbClient)
	validation.AssertNotNil("logger", logger)

	return &gooseMigratorAdapter{
		name:     name,
		cfg:      cfg,
		dbClient: dbClient,
		logger:   logger,
	}
}

func (g *gooseMigratorAdapter) Down(version uint) error {
	panic("unimplemented")
}

func (g *gooseMigratorAdapter) Up(version uint) error {
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
