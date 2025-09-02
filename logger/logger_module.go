package logger

import (
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/logger/config"
)

func ModuleFunc(loggerCfg config.Config) module.Module {
	m := module.ModuleFunc("logger")

	m.AddProviders(provider.ProvideFunc(func() config.Config { return loggerCfg }))

	return m
}
