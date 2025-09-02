package database_client

import (
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/database/database_client/config"
)

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"database_client",

		provider.ProvideFunc(config.LoadEnvConfig),
	)
}
