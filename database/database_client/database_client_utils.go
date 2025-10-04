package databaseclient

import (
	"frisboo-bank/pkg/database/database_client/adapters/mongodb"
	"frisboo-bank/pkg/database/database_client/adapters/postgres"
	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/database/database_client/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
)

func NoDatabaseClientOfTypeError(name string, sType databaseclienttype.DatabaseClientType) error {
	return syserrors.Newf("database-client type:{%s} for server {%s} does not exist", sType, name)
}

func GetInstance(name string, cfg *config.Config, logger loggerContracts.Logger) (contracts.DatabaseClient, error) {
	var adapter contracts.DatabaseClient
	var err error

	switch cfg.Type {
	case databaseclienttype.DatabaseClientTypes.MONGODB:
		adapter, err = mongodb.New(name, cfg, logger)
	case databaseclienttype.DatabaseClientTypes.POSTGRES:
		adapter, err = postgres.New(name, cfg, logger)
	default:
		return nil, NoDatabaseClientOfTypeError(name, cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return New(adapter), nil
}
