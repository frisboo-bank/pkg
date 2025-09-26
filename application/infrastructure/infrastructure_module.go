package infrastructure

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	httpServer "frisboo-bank/pkg/http/http_server"
	httpServerConfig "frisboo-bank/pkg/http/http_server/config"
	rpcServer "frisboo-bank/pkg/rpc/rpc_server"
	rpcServerConfig "frisboo-bank/pkg/rpc/rpc_server/config"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	m := module.ModuleFunc("infrastructure")

	httpServerModule, err := registerHTTPServerModule(appBuilder)
	if err != nil {
		appBuilder.Logger().Panic(err)
	}
	m.AddModule(httpServerModule)

	rpcServerModule, err := registerRPCServerModule(appBuilder)
	if err != nil {
		appBuilder.Logger().Panic(err)
	}
	m.AddModule(rpcServerModule)

	// m.AddModule(health.ModuleFunc())

	return m
}

func registerHTTPServerModule(appBuilder applicationContracts.ApplicationBuilder) (module.Module, error) {
	reg, err := httpServerConfig.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
	if err != nil {
		return nil, err
	}
	return httpServer.ModuleFunc(reg), nil
}

func registerRPCServerModule(appBuilder applicationContracts.ApplicationBuilder) (module.Module, error) {
	reg, err := rpcServerConfig.LoadRegistry(appBuilder.ConfigLoader(), appBuilder.Environment())
	if err != nil {
		return nil, err
	}
	return rpcServer.ModuleFunc(reg), nil
}
