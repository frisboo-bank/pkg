package rpcserver

import (
	"fmt"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"
	"frisboo-bank/pkg/rpc/rpc_server/grpc"
)

func GetInstance(
	sType rpcservertype.RpcServerType,
	logger loggerContracts.Logger,
	opt *options.OptionBuilder[config.Config],
) (contracts.RPCServer, error) {
	var adapter contracts.RPCServerAdapter

	switch sType {
	case rpcservertype.RpcServerTypes.GRPC:
		adapter = grpc.New(logger)
	default:
		return nil, fmt.Errorf("(rpc-server-factory) no server of type `%q` exists", sType)
	}

	server, err := New(adapter, logger, opt)
	if err != nil {
		return nil, err
	}

	return server, nil
}
