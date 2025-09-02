package infrastructure

import (
	appConfig "frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/health"
	httpserver "frisboo-bank/pkg/http/http_server"
)

func ModuleFunc(cfg *appConfig.Config) module.Module {
	// if cfg {
	// 	deps = append(deps, httpserver.Module)
	// }
	//
	// if cfg.EnableGRPCServer {
	// 	deps = append(deps, rpcserver.Module)
	// }

	m := module.ModuleFunc(
		"infrastructure",

		httpserver.ModuleFunc(),

		health.ModuleFunc(),
	)

	return m
}
