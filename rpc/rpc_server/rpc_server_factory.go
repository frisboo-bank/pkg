package rpcserver

import (
	"frisboo-bank/pkg/rpc/rpc_server/adapters/grpc"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/syserrors"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
)

func NoServerOfTypeError(sType rpcservertype.RpcServerType) error {
	return syserrors.Newf("no server of type `%q` exists", sType)
}

func GetInstance(cfg *config.Config) (contracts.RPCServer, error) {
	var adapter contracts.RPCServerAdapter

	switch cfg.Type {
	case rpcservertype.RpcServerTypes.GRPC:
		adapter = grpc.New(cfg, nil)
	default:
		return nil, NoServerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
