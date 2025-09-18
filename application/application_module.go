package application

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/application/infrastructure"
	"frisboo-bank/pkg/container/dependencies/module"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	return module.ModuleFunc("application",
		infrastructure.ModuleFunc(appBuilder),
	)
}
