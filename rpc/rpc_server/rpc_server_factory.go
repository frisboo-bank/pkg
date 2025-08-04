package rpcserver

import (
	"fmt"

	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/grpc"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

func GetInstanceFromConfig(config *config.RPCServerConfig, logger loggerContracts.Logger) (contracts.RPCServer, error) {
	instance, err := GetInstance(config.Type, logger)
	if err != nil {
		return nil, err
	}

	return instance.WithConfig(config), nil
}

func GetInstance(rpcServerType rpcservertype.RpcServerType, logger loggerContracts.Logger) (contracts.RPCServer, error) {
	switch rpcServerType {
	case rpcservertype.RpcServerTypes.GRPC:
		return grpc.New(logger), nil
	default:
		return nil, fmt.Errorf("(rpc-server-factory) no server of type `%q` exists", rpcServerType)
	}
}
