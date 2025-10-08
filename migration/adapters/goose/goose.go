package goose

import (
	"context"

	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	vGoose "github.com/pressly/goose/v3"
)

var _ contracts.MigratorAdapter = (*gooseMigratorAdapter)(nil)

type gooseMigratorAdapter struct {
	name     string
	cfg      *config.Config
	ctx      context.Context
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

	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}

	return &gooseMigratorAdapter{
		name:     name,
		cfg:      cfg,
		ctx:      ctx,
		dbClient: dbClient,
		logger:   logger,
	}
}

func (g *gooseMigratorAdapter) Down(version uint) error {
	return g.Run(migrationcommandtype.MigrationCommandTypes.DOWN, version)
}

func (g *gooseMigratorAdapter) Up(version uint) error {
	return g.Run(migrationcommandtype.MigrationCommandTypes.UP, version)
}

func (g *gooseMigratorAdapter) Reset() error {
	return vGoose.ResetContext(g.ctx, g.dbClient.DB(), g.cfg.MigrationsDir)
}

func (g *gooseMigratorAdapter) Run(command migrationcommandtype.MigrationCommandType, version uint) error {
	switch command {
	case migrationcommandtype.MigrationCommandTypes.UP:
		if version == 0 {
			return vGoose.RunContext(g.ctx, "up", g.dbClient.DB(), g.cfg.MigrationsDir)
		}
		return vGoose.RunContext(g.ctx, "up-to VERSION ", g.dbClient.DB(), g.cfg.MigrationsDir)
	case migrationcommandtype.MigrationCommandTypes.DOWN:
		if version == 0 {
			return vGoose.RunContext(g.ctx, "down", g.dbClient.DB(), g.cfg.MigrationsDir)
		}
		return vGoose.RunContext(g.ctx, "down-to VERSION ", g.dbClient.DB(), g.cfg.MigrationsDir)
	default:
		return syserrors.Newf("invalid command:{%s}", command)
	}
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
