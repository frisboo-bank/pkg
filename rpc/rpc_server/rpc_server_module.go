package rpcserver

import (
	"frisboo-bank/pkg/container/dependencies/module"
)

func ModuleFunc() module.Module {
	m := module.ModuleFunc("rpc-server")

	// m.AddHooks(
	// 	hook.HooksFunc(func(rpcServer contracts.RPCServer) waiterContracts.WaitFunc {
	// 		return func(ctx context.Context) error {
	// 			var err error
	//
	// 			go func() {
	// 				err = rpcServer.Start(ctx)
	// 			}()
	// 			return err
	// 		}
	// 	}, func(rpcServer contracts.RPCServer) waiterContracts.CleanupFunc {
	// 		return func(ctx context.Context) error {
	// 			return rpcServer.Shutdown(ctx)
	// 		}
	// 	}),
	// )

	return m
}
