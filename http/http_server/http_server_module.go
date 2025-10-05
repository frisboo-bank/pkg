package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	applicationContracts "frisboo-bank/pkg/application/contracts"

	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

const HTTPServersGroup = "http-servers"

func ModuleFunc(appBuilder applicationContracts.ApplicationBuilder) module.Module {
	validation.AssertNotNil("appBuilder", appBuilder)

	configLoader := appBuilder.ConfigLoader()
	env := appBuilder.Environment()
	logger := appBuilder.Logger()

	// Load and register the config registry
	cfgRegistry, err := config.LoadRegistry(configLoader, env)
	if err != nil {
		logger.Fatalf("failed to register http-server module with error: %v", err)
	}

	m := module.ModuleFunc(
		"http-server",
		provider.ProvideFunc(func() config.Registry { return cfgRegistry }),
	)

	for _, name := range cfgRegistry.Names() {
		cfg, err := cfgRegistry.GetByName(name)
		if err != nil {
			logger.Fatalf("failed to register http-server:{%s} module with error:{%v}", name, err)
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

	log.Debugf("Try to register http-server:{%s} module", name)

	m := module.ModuleFunc("http-server:" + name)

	// Instance registration name
	providerName := "http-server:" + name

	m.AddProvider(provider.ProvideFunc(func(loggerCfgRegistry loggerConfig.Registry, appLogger loggerContracts.Logger) (contracts.HTTPServer, error) {
		// Resolve logger (either server-specific or fallback to app logger)
		log, err := logger.GetByNameWithFallback(loggerCfgRegistry, cfg.Logger, appLogger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "http-server:{%s} logger", name)
		}
		return GetInstance(name, cfg, log)
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

				srv.Logger().Infof("http-server:{%s} is listening on address:{%s}", srv.Name(), srv.Config().Address())

				if err := srv.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
					srv.Logger().Fatalf("{%s} failed to start with error:{%v}", srv.Name(), err)
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
					srv.Logger().Errorf("http-server:{%s} shutdown failed with error:{%v}", srv.Name(), err)
				} else {
					srv.Logger().Infof("http-server:{%s} shutdown successfully", srv.Name())
				}
				return nil
			}
		},
		hook.NamedDep("httpServerRef", providerName),
	))

	return m
}
