package rpcserver

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

var Module = module.ModuleFunc(
	"rpc-server",

	provider.ProvideFunc(config.Load),

	// provider.ProvideFunc(
	// 	func(loggerEnvCfg *loggerConfig.EnvConfig, envCfg *config.EnvConfig) (contracts.RPCServer, error) {
	// 		loggerOpts := loggerConfig.FromEnvConfig(loggerEnvCfg).
	// 			With(loggerConfig.Prefix("rpc-server"))
	//
	// 		logger, err := logger.GetInstance(loggerEnvCfg.Type, loggerOpts)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		opts := config.FromEnvConfig(envCfg)
	//
	// 		rpcServer, err := GetInstance(envCfg.Type, logger, opts)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		return rpcServer, nil
	// 	},
	// ),
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
