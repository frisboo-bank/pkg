package rpcserver

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = module.ModuleFunc(
	"rpc-server",

	//
	hook.HooksFunc(func(rpcServer contracts.RPCServer) waiterContracts.WaitFunc {
		return func(ctx context.Context) error {
			var err error

			go func() {
				err = rpcServer.Start(ctx)
			}()
			return err
		}
	}, func(rpcServer contracts.RPCServer) waiterContracts.CleanupFunc {
		return func(ctx context.Context) error {
			return rpcServer.Shutdown(ctx)
		}
	}),
)
