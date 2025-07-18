package database_client

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/database/database_client/options"
	"frisboo-bank/pkg/environment"
)

var Module = container.NewModule(
	"database_client",

	container.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.DatabaseClientOptions, error) {
			return options.ProvideDatabaseClientOptions(loader, env)
		},
	),
)
