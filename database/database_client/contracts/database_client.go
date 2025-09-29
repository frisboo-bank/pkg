package contracts

import (
	"frisboo-bank/pkg/database/database_client/config"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	databaseClientCommon interface {
		Ping() error
		Disconnect() error
		Name() string
		Type() databaseclienttype.DatabaseClientType
		Config() *config.Config
		Logger() loggerContracts.Logger
	}

	DatabaseClient interface {
		databaseClientCommon
	}

	DatabaseClientAdapter interface {
		databaseClientCommon
	}
)
