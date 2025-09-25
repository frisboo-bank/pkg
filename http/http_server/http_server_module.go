package httpserver

import (
	"context"
	"errors"
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

const HTTPServersGroup = "http_servers"

func FailedToStartError(err error) error { return syserrors.Wrap(err, "server failed to start") }
func FailedToStopError(err error) error  { return syserrors.Wrap(err, "server failed to stop") }

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
		m.AddModules(serverModuleFunc(name, &cfg))
	}

	return m
}

func serverModuleFunc(instance string, cfg *config.Config) module.Module {
	validation.AssertNotEmpty("instance", instance)

	m := module.ModuleFunc("http_server:" + instance)

	// Instance registration name
	providerName := "http_server:" + instance

	m.AddProviders(provider.ProvideFunc(func(
		loggerCfgRegistry loggerConfig.Registry,
		appLogger loggerContracts.Logger,
	) (contracts.HTTPServer, error) {
		// Resolve logger (either server-specific or fallback to app logger)
		log := appLogger
		if cfg.Logger != "" {
			loggerCfg, err := loggerCfgRegistry.GetByName(cfg.Logger)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to load http-server %s logger config %s", instance, cfg.Logger)
			}
			log, err = logger.GetInstance(&loggerCfg)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to initialize http-server %s logger %s", instance, cfg.Logger)
			}
		}

		// Select proper adapter
		var adapter contracts.HTTPServerAdapter
		switch cfg.Type {
		case httpservertype.HttpServerTypes.ECHO:
			adapter = echo.New(cfg, log, nil)
		default:
			return nil, syserrors.Newf("http_server %s is using an invalid type: got %s", instance, cfg.Type)
		}

		return New(adapter), nil
	},
		provider.Name(providerName),
		provider.Group(HTTPServersGroup),
	))

	type hookParams struct {
		HTTPServer contracts.HTTPServer `name:"httpServerRef"`
	}

	m.AddHooks(
		hook.HooksFunc(
			func(p hookParams) waiterContracts.WaitFunc {
				return func(ctx context.Context) error {
					srv := p.HTTPServer
					srv.SetupDefaultMiddlewares()

					srv.Logger().Printf("http server started; listening at %s", srv.Config().Address())
					defer srv.Logger().Print("http server shutdown")

					if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
						return FailedToStartError(err)
					}
					return nil
				}
			},
			func(p hookParams) waiterContracts.CleanupFunc {
				return func(ctx context.Context) error {
					srv := p.HTTPServer

					srv.Logger().Print("http server stopping")

					ctx, cancel := context.WithTimeout(ctx, srv.Config().ServerShutdownTimeout)
					defer cancel()

					if err := srv.Stop(ctx); err != nil {
						return FailedToStopError(err)
					}
					return nil
				}
			},
			hook.NamedDep("httpServerRef", providerName),
		),
	)

	return m
}
