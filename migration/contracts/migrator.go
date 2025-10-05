package contracts

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
)

type migratorCommon interface {
	Up(version uint) error
	Down(version uint) error
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
