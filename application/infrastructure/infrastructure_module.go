package infrastructure

import (
	"frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/container/dependencies/module"
	httpserver "frisboo-bank/pkg/http/http_server"
	rpcserver "frisboo-bank/pkg/rpc/rpc_server"
)

func ModuleFunc(appCfg *config.Config) module.Module {
	m := module.ModuleFunc("infrastructure",
		httpserver.ModuleFunc(),
		rpcserver.ModuleFunc(),
	)

	m.AddModules(
	// health.ModuleFunc(),
	)

	return m
}
