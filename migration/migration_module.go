package migration

import (
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/migration/config"
)

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"migration",

		provider.ProvideFunc(config.LoadEnvConfig),
	)
}
