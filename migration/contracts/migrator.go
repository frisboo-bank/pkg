package contracts

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
)

type migratorCommon interface {
	Up(version uint) error
	Down(version uint) error
	Reset() error
	Run(command migrationcommandtype.MigrationCommandType, version uint) error

	Name() string
	Type() migratortype.MigratorType
	Config() *config.Config
	Logger() loggerContracts.Logger
}

type Migrator interface {
	migratorCommon
}

type MigratorAdapter interface {
	migratorCommon
}
