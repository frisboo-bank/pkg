package migration

import (
	"frisboo-bank/pkg/container/dependencies/module"
)

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"migration",
	)
}
