package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/http/http_server/adapters/echo"
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

const HTTPServersGroup = "http-servers"

func ModuleFunc(registry config.Registry) module.Module {
	m := module.ModuleFunc("http-server")

	for _, name := range registry.Names() {
		cfg, err := registry.GetByName(name)
		if err != nil {
			panic(syserrors.Wrapf(err, "failed to load config for http-server %s", name))
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
	providerName := "http-server:" + name

	m.AddProvider(provider.ProvideFunc(func(
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
			log, err = logger.GetInstance(&loggerCfg)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to initialize http-server %s logger %s", name, cfg.Logger)
			}
		}

		// Select proper adapter
		var adapter contracts.HTTPServerAdapter
		switch cfg.Type {
		case httpservertype.HttpServerTypes.ECHO:
			adapter = echo.New(name, cfg, log, nil)
		default:
			return nil, syserrors.Newf("http-server %s is using an invalid type: got %s", name, cfg.Type)
		}

		return New(adapter), nil
	},
		provider.Name(providerName),
		provider.Group(HTTPServersGroup),
	))

	type hookParams struct {
		HTTPServer contracts.HTTPServer `name:"httpServerRef"`
	}

	m.AddHook(hook.HooksFunc(fmt.Sprintf("http-server-%s-hook", name),
		func(p hookParams) waiterContracts.WaitFunc {
			return func(ctx context.Context) error {
				srv := p.HTTPServer
				srv.SetupDefaultMiddlewares()

				srv.Logger().Infof("%s is listening on address:{%s}", srv.Name(), srv.Config().Address())

				if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					srv.Logger().Fatalf("%s failed to start with error: %v", srv.Name(), err)
				}

				return nil
			}
		},
		func(p hookParams) waiterContracts.CleanupFunc {
			return func(ctx context.Context) error {
				srv := p.HTTPServer

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
		hook.NamedDep("httpServerRef", providerName),
	))

	return m
}
