package httpserver

import (
	"context"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/http/http_server/adapters/gin"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

const HTTPServersGroup = "http_servers"

func ServerFailedToStartError(err error) error {
	return syserrors.Wrap(err, "server failed to start with error", "HTTPServer", "Hooks", "Start")
}

func ServerFailedToStopError(err error) error {
	return syserrors.Wrap(err, "server failed to stop", "HTTPServer", "Hooks", "Stop")
}

func ModuleFunc() module.Module {
	m := module.ModuleFunc("http_server")

	m.AddProviders(
		provider.ProvideFunc(
			func(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (config.Registry, error) {
				return config.LoadRegistry(configLoader, env)
			},
		),
	)

	m.AddProviders(provider.ProvideFunc(func(
		cfgReg config.Registry,
		loggerCfgReg loggerConfig.Registry,
		appLogger loggerContracts.Logger,
	) ([]contracts.HTTPServer, error) {
		names := cfgReg.Names()
		servers := make([]contracts.HTTPServer, 0, len(names))

		for _, name := range cfgReg.Names() {
			cfg, err := cfgReg.GetByName(name)
			if err != nil {
				return nil, syserrors.Wrapf(err, "failed to load http-server %s config", name)
			}

			// Skip disabled servers
			if !cfg.Enabled {
				continue
			}

			// Resolve logger (either server-specific or fallback to app logger)
			srvLogger := appLogger
			if cfg.Logger != "" {
				loggerCfg, err := loggerCfgReg.GetByName(cfg.Logger)
				if err != nil {
					return nil, syserrors.Wrapf(err, "failed to load http-server %s logger config %s", name, cfg.Logger)
				}
				srvLogger, err = logger.GetInstance(&loggerCfg)
				if err != nil {
					return nil, syserrors.Wrapf(err, "failed to initialize http-server %s logger %s", name, cfg.Logger)
				}
			}

			// Select proper adapter
			var adapter contracts.HTTPServerAdapter
			switch cfg.Type {
			case httpservertype.HttpServerTypes.GIN:
				adapter = gin.New(&cfg, srvLogger)
			default:
				return nil, syserrors.Newf("http_server %s is using an invalid type: got %s", name, cfg.Type)
			}
			server := New(adapter)

			servers = append(servers, server)
		}

		for _, server := range servers {
			s := server
			m.AddHooks(hook.HooksFunc(func() waiterContracts.WaitFunc {
				return func(ctx context.Context) error {
					s.SetupDefaultMiddlewares()
					errChan := make(chan error, 1)
					go func() {
						errChan <- s.Start(ctx)
					}()
					select {
					case err := <-errChan:
						if err != nil {
							return ServerFailedToStartError(err)
						}
						return nil
					case <-ctx.Done():
						// Context cancelled, attempt to stop the server
						stopErr := s.Stop(ctx)
						// Wait for the server goroutine to finish
						err := <-errChan
						if err != nil {
							return ServerFailedToStartError(err)
						}
						if stopErr != nil {
							return ServerFailedToStopError(stopErr)
						}
						return ctx.Err()
					}
				}
			}, func() waiterContracts.CleanupFunc {
				return func(ctx context.Context) error {
					if err := s.Stop(ctx); err != nil {
						return ServerFailedToStopError(err)
					}
					return nil
				}
			}))
		}

		return servers, nil
	}))

	return m
}
