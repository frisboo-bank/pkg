package application

import (
	"frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/application/infrastructure"
	"frisboo-bank/pkg/container/dependencies/module"
)

func ModuleFunc(appCfg *config.Config) module.Module {
	return module.ModuleFunc("application",
		infrastructure.ModuleFunc(appCfg),
	)
}
