package database_client

import (
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/database/database_client/options"
	"frisboo-bank/pkg/environment"
)

var Module = module.NewModule(
	"database_client",

	provider.Provide(
		func(loader configContracts.ConfigLoader, env environment.Environment) (*options.DatabaseClientOptions, error) {
			return options.ProvideDatabaseClientOptions(loader, env)
		},
	),
)
