package infrastructure

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	databaseclient "frisboo-bank/pkg/database/database_client"
	httpServer "frisboo-bank/pkg/http/http_server"
	rpcServer "frisboo-bank/pkg/rpc/rpc_server"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	m := module.ModuleFunc("infrastructure",
		httpServer.ModuleFunc(appBuilder),
		rpcServer.ModuleFunc(appBuilder),
		databaseclient.ModuleFunc(appBuilder),
	)

	// m.AddModule(health.ModuleFunc())

	return m
}
