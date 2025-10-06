package rpcserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	applicationContracts "frisboo-bank/pkg/application/contracts"
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

const (
	RPCServersGroup    = "rpc-servers"
	RPCServersProvider = "rpc-server:%s"
)

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	validation.AssertNotNil("appBuilder", appBuilder)

	configLoader := appBuilder.ConfigLoader()
	env := appBuilder.Environment()
	logger := appBuilder.Logger()

	// Load and register the config registry
	cfgRegistry, err := config.LoadRegistry(configLoader, env)
	if err != nil {
		logger.Panicw(
			"failed to register rpc-server module",
			loggerContracts.Fields{"err": err, "cause": syserrors.Cause(err)},
		)
	}

	m := module.ModuleFunc(
		"rpc-server",
		provider.ProvideFunc(func() config.Registry { return cfgRegistry }),
	)

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			logger.Panicw(
				"failed to register rpc-server module",
				loggerContracts.Fields{"err": err, "cause": syserrors.Cause(err)},
			)
		}
		if !cfg.Enabled {
			continue
		}
		m.AddModule(serverModuleFunc(name, logger, &cfg))
	}

	return m
}

func serverModuleFunc(name string, log loggerContracts.Logger, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("log", log)
	validation.AssertNotNil("cfg", cfg)

	log.Debugf("Try to register rpc-server:{%s} module", name)

	m := module.ModuleFunc("rpc-server:" + name)

	type providerProps struct {
		LoggerCfgRegistry loggerConfig.Registry
		AppLogger         loggerContracts.Logger
	}

	m.AddProvider(provider.ProvideFunc(func(props providerProps) (contracts.RPCServer, error) {
		loggerCfgRegistry := props.LoggerCfgRegistry
		appLogger := props.AppLogger

		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "rpc-server %s logger", name)
		}

		return GetInstance(name, cfg, log)
	},
		provider.Name(fmt.Sprintf(RPCServersProvider, name)),
		provider.Group(RPCServersGroup),
	))

	type hookProps struct {
		RPCServer contracts.RPCServer `name:"rpcServerRef"`
	}

	m.AddHook(hook.HooksFunc(fmt.Sprintf("rpc-server-%s-hook", name),
		func(props hookProps) waiterContracts.WaitFunc {
			return func(ctx context.Context) error {
				srv := props.RPCServer
				srv.Logger().Infof("%s is listening on address:{%s}", srv.Name(), srv.Config().Address())

				if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					srv.Logger().Fatalf("%s failed to start with error: %v", srv.Name(), err)
				}

				return nil
			}
		},
		func(props hookProps) waiterContracts.CleanupFunc {
			return func(ctx context.Context) error {
				srv := props.RPCServer

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
		hook.NamedDep("rpcServerRef", fmt.Sprintf(RPCServersProvider, name)),
	))

	return m
}
