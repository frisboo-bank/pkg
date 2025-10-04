package migration

import (
	databaseclientContracts "frisboo-bank/pkg/database/database_client/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/migration/adapters/goose"
	"frisboo-bank/pkg/migration/config"
	"frisboo-bank/pkg/migration/contracts"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/syserrors"
)

func NoMigratorOfTypeError(name string, sType migratortype.MigratorType) error {
	return syserrors.Newf("migrator type:{%s} for client {%s} does not exist", sType, name)
}

func GetInstance(name string, cfg *config.Config, dbClient databaseclientContracts.DatabaseClient, logger loggerContracts.Logger) (contracts.Migrator, error) {
	var adapter contracts.Migrator

	switch cfg.Type {
	case migratortype.MigratorTypes.GOOSE:
		adapter = goose.New(name, cfg, dbClient, logger)
	default:
		return nil, NoMigratorOfTypeError(name, cfg.Type)
	}

	return New(adapter), nil
}
