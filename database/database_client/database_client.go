package databaseclient

import (
	"database/sql"

	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/database/database_client/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"
)

var _ contracts.DatabaseClient = (*databaseClient)(nil)

type databaseClient struct {
	adapter contracts.DatabaseClientAdapter
}

func New(adapter contracts.DatabaseClientAdapter) contracts.DatabaseClient {
	validation.AssertNotNil("adapter", adapter)
	return &databaseClient{adapter}
}

func (d *databaseClient) Ping() error {
	return d.adapter.Ping()
}

func (d *databaseClient) Disconnect() error {
	return d.adapter.Disconnect()
}

func (d *databaseClient) Name() string {
	return d.adapter.Name()
}

func (d *databaseClient) Type() databaseclienttype.DatabaseClientType {
	return d.adapter.Type()
}

func (d *databaseClient) Config() *config.Config {
	return d.adapter.Config()
}

func (d *databaseClient) DB() *sql.DB {
	return d.adapter.DB()
}

func (d *databaseClient) Logger() loggerContracts.Logger {
	return d.adapter.Logger()
}
