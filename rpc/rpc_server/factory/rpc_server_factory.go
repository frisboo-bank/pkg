package factory

import (
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/grpc"
	"frisboo-bank/pkg/rpc/rpc_server/options"
)

func GetInstance(config *options.RPCServerOptions, configs ...options.RPCServerOption) (contracts.RPCServer, error) {
	for _, c := range configs {
		c(config)
	}

	return grpc.NewGRPCServer(config), nil
}
