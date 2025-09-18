package httpserver

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/http/http_server/adapters/gin"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

const HTTPServersGroup = "http_servers"

func FailedToStartError(err error) error {
	return syserrors.Wrap(err, "server failed to start with error", "HTTPServer", "Hooks", "Start")
}

func FailedToStopError(err error) error {
	return syserrors.Wrap(err, "server failed to stop", "HTTPServer", "Hooks", "Stop")
}

func ModuleFunc(reg config.Registry) module.Module {
	m := module.ModuleFunc("http_server")

	for _, name := range reg.Names() {
		cfg, err := reg.GetByName(name)
		if err != nil {
			panic(syserrors.Wrapf(err, "failed to load config for http-server %s", name))
		}
		if !cfg.Enabled {
			continue
		}
		m.AddModules(serverModuleFunc(name, cfg))
	}

	return m
}

func serverModuleFunc(name string, cfg config.Config) module.Module {
	validation.AssertNotEmpty("name", name)

	m := module.ModuleFunc("http_server:" + name)

	m.AddProviders(provider.ProvideFunc(func(
		loggerCfgRegistry loggerConfig.Registry,
		appLogger loggerContracts.Logger,
	) (contracts.HTTPServer, error) {
		// Resolve logger (either server-specific or fallback to app logger)
		log := appLogger
		if cfg.Logger != "" {
			loggerCfg, err := loggerCfgRegistry.GetByName(cfg.Logger)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to load http-server %s logger config %s", name, cfg.Logger)
			}
			log, err = logger.GetInstance(loggerCfg)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to initialize http-server %s logger %s", name, cfg.Logger)
			}
		}

		// Select proper adapter
		var adapter contracts.HTTPServerAdapter
		switch cfg.Type {
		case httpservertype.HttpServerTypes.GIN:
			adapter = gin.New(&cfg, log)
		default:
			return nil, syserrors.Newf("http_server %s is using an invalid type: got %s", name, cfg.Type)
		}

		return New(adapter), nil
	}))

	m.AddProviders(
		provider.ProvideFunc(func(s contracts.HTTPServer) contracts.HTTPServer { return s }, provider.Name("http_server:"+name)),
		provider.ProvideFunc(func(s contracts.HTTPServer) contracts.HTTPServer { return s }, provider.Group("http_servers")),
	)

	m.AddHooks(hook.HooksFunc(httpServerHookStart, httpServerHookStop))

	return m
}

func httpServerHookStart(srv contracts.HTTPServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		srv.SetupDefaultMiddlewares()

		errChan := make(chan error, 1)
		go func() {
			errChan <- srv.Start(ctx)
		}()

		select {
		case err := <-errChan:
			if err != nil {
				return FailedToStartError(err)
			}
			return nil
		case <-ctx.Done():
			if err := srv.Stop(ctx); err != nil {
				return FailedToStopError(err)
			}
			return ctx.Err()
		}
	}
}

func httpServerHookStop(srv contracts.HTTPServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		if err := srv.Stop(ctx); err != nil {
			return FailedToStopError(err)
		}
		return nil
	}
}
