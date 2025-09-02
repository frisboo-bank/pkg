package application

import (
	"frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
)

func ModuleFunc(cfg *config.Config) module.Module {
	return module.ModuleFunc(
		"application",

		provider.ProvideFunc(func() *config.Config { return cfg }),
	)
}
