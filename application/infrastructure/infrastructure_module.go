package infrastructure

import (
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/health"
	httpserver "frisboo-bank/pkg/http/http_server"
	rpcserver "frisboo-bank/pkg/rpc/rpc_server"
)

var Module = container.NewModule("infrastructure",
	httpserver.Module,
	rpcserver.Module,
	health.Module,
)
