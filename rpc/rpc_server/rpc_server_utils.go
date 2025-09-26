package rpcserver

import (
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/adapters/grpc"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
	"frisboo-bank/pkg/syserrors"
)

func NoRPCServerOfTypeError(name string, sType rpcservertype.RpcServerType) error {
	return syserrors.Newf("rpc-server type %s for server %s does not exist", sType, name)
}

func GetInstance(name string, cfg *config.Config, logger loggerContracts.Logger) (contracts.RPCServer, error) {
	var adapter contracts.RPCServerAdapter

	switch cfg.Type {
	case rpcservertype.RpcServerTypes.GRPC:
		adapter = grpc.New(name, cfg, logger, nil)
	default:
		return nil, NoRPCServerOfTypeError(name, cfg.Type)
	}

	return New(adapter), nil
}
