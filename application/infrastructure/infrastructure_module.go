package infrastructure

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	httpServer "frisboo-bank/pkg/http/http_server"
	httpServerConfig "frisboo-bank/pkg/http/http_server/config"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	m := module.ModuleFunc("infrastructure")

	httpServerModule, err := registerHTTPServerModule(appBuilder)
	if err != nil {
		panic(err)
	}
	m.AddModules(httpServerModule)

	m.AddModules(
	// health.ModuleFunc(),
	)

	return m
}

func registerHTTPServerModule(appBuilder applicationContracts.ApplicationBuilder) (module.Module, error) {
	reg, err := httpServerConfig.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
	if err != nil {
		return nil, err
	}
	return httpServer.ModuleFunc(reg), nil
}
