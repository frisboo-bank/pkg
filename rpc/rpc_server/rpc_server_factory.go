package rpcserver

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/grpc"
)

func GetInstance(logger loggerContracts.Logger) (contracts.RPCServer, error) {
	return grpc.NewGRPCServer(logger), nil
}
