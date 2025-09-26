package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/config"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

const RPCServersGroup = "rpc-servers"

func ModuleFunc(registry config.Registry) module.Module {
	m := module.ModuleFunc("rpc-server")

	for _, name := range registry.Names() {
		cfg, err := registry.GetByName(name)
		if err != nil {
			panic(syserrors.Wrapf(err, "failed to load config for rpc-server %s", name))
		}
		if !cfg.Enabled {
			continue
		}
		m.AddModule(serverModuleFunc(name, &cfg))
	}

	return m
}

func serverModuleFunc(name string, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("name", name)

	m := module.ModuleFunc("http-server:" + name)

	// Instance registration name
	providerName := "rpc-server:" + name

	m.AddProvider(provider.ProvideFunc(func(
		loggerCfgRegistry loggerConfig.Registry,
		appLogger loggerContracts.Logger,
	) (contracts.RPCServer, error) {
		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "rpc-server %s logger", name)
		}

		return GetInstance(name, cfg, log)
	},
		provider.Name(providerName),
		provider.Group(RPCServersGroup),
	))

	type hookParams struct {
		RPCServer contracts.RPCServer `name:"rpcServerRef"`
	}

	m.AddHook(hook.HooksFunc(fmt.Sprintf("rpc-server-%s-hook", name),
		func(p hookParams) waiterContracts.WaitFunc {
			return func(ctx context.Context) error {
				srv := p.RPCServer
				srv.Logger().Infof("%s is listening on address:{%s}", srv.Name(), srv.Config().Address())

				if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					srv.Logger().Fatalf("%s failed to start with error: %v", srv.Name(), err)
				}

				return nil
			}
		},
		func(p hookParams) waiterContracts.CleanupFunc {
			return func(ctx context.Context) error {
				srv := p.RPCServer

				ctx, cancel := context.WithTimeout(ctx, srv.Config().ServerShutdownTimeout)
				defer cancel()

				if err := srv.Stop(ctx); err != nil {
					srv.Logger().Errorf("shutting down %s failed with error: %v", srv.Name(), err)
				} else {
					srv.Logger().Infof("%s server shutdown successfully", srv.Name())
				}
				return nil
			}
		},
		hook.NamedDep("rpcServerRef", providerName),
	))

	return m
}
