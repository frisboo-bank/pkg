package rpcserver

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/grpc"
	"frisboo-bank/pkg/rpc/rpc_server/options"
)

func GetInstanceFromOptions(
	options *options.RPCServerOptions,
	logger loggerContracts.Logger,
) (contracts.RPCServer, error) {
	instance, err := GetInstance(logger)
	if err != nil {
		return nil, err
	}

	return instance.WithOptions(options), nil
}

func GetInstance(logger loggerContracts.Logger) (contracts.RPCServer, error) {
	return grpc.NewGRPCServer(logger), nil
}
