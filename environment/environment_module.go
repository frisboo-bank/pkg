package environment

import (
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
)

func ModuleFunc(env Environment) module.Module {
	return module.ModuleFunc(
		"environment",
		provider.ProvideFunc(func() Environment { return env }),
	)
}
