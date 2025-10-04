package migration

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/validation"
)

var _ contracts.Migrator = (*migrator)(nil)

type migrator struct {
	adapter contracts.MigratorAdapter
}

func New(adapter contracts.MigratorAdapter) contracts.Migrator {
	validation.AssertNotNil("adapter", adapter)
	return &migrator{adapter}
}

func (m *migrator) Run(commandType migrationcommandtype.MigrationCommandType, version string) error {
	return m.adapter.Run(commandType, version)
}

func (m *migrator) Config() *config.Config {
	return m.adapter.Config()
}

func (m *migrator) Logger() loggerContracts.Logger {
	return m.adapter.Logger()
}

func (m *migrator) Name() string {
	return m.adapter.Name()
}

func (m *migrator) Type() migratortype.MigratorType {
	return m.adapter.Type()
}
