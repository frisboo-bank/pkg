package database_client

import (
	"frisboo-bank/pkg/container/dependencies/module"
)

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"database_client",
	)
}
