package application

import (
	appConfig "frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/application/contracts"
	configloader "frisboo-bank/pkg/config/config_loader"
	configloaderConfig "frisboo-bank/pkg/config/config_loader/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/container/adapters/dig"
	containerConfig "frisboo-bank/pkg/container/config"
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	containerEnums "frisboo-bank/pkg/container/enums"
	"frisboo-bank/pkg/environment"
	httpServerEnums "frisboo-bank/pkg/http/http_server/enums"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerEnums "frisboo-bank/pkg/logger/enums"
	rpcServerEnums "frisboo-bank/pkg/rpc/rpc_server/enums"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/waiter"
	waiterConfig "frisboo-bank/pkg/waiter/config"
)

var _ contracts.ApplicationBuilder = (*applicationBuilder)(nil)

type applicationBuilder struct {
	environment          environment.Environment
	logger               loggerContracts.Logger
	configLoader         configloaderContracts.ConfigLoader
	loggerConfigRegistry loggerConfig.Registry
	appConfig            appConfig.Config
	container            containerContracts.Container
	modules              []module.Module
	providers            []provider.Provider
	decorators           []decorator.Decorator
}

func NewApplicationBuilder(environments ...environment.Environment) (contracts.ApplicationBuilder, error) {
	env, err := environment.Load(environments...)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to instantiate environment")
	}

	configLoader, err := configloader.New(
		configloaderConfig.Debug(false),
		configloaderConfig.DecodeHookFuncs(
			containerEnums.ContainerEnumsDecodeHook(),
			environment.EnvironmentEnumsDecodeHook(),
			httpServerEnums.HTTPServerEnumsDecodeHook(),
			loggerEnums.LoggerEnumsDecodeHook(),
			rpcServerEnums.RPCServerEnumsDecodeHook(),
		),
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to instantiate config-loader")
	}

	loggerCfgRegistry, err := loggerConfig.LoadRegistry(configLoader, env)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to initialize logger config registry")
	}

	appCfg, err := appConfig.Load(configLoader, env)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load app config")
	}

	// Initialize logger
	var appLoggerCfg loggerConfig.Config
	if appCfg.Logger != "" {
		appLoggerCfg, err = loggerCfgRegistry.GetByName(appCfg.Logger)
	} else {
		appLoggerCfg, err = loggerCfgRegistry.GetDefault()
	}
	if err != nil {
		return nil, err
	}

	appLogger, err := logger.GetInstance(&appLoggerCfg)
	if err != nil {
		return nil, err
	}

	appBuilder := &applicationBuilder{
		environment:          env,
		logger:               appLogger,
		configLoader:         configLoader,
		loggerConfigRegistry: loggerCfgRegistry,
		appConfig:            appCfg,
	}

	if err := appBuilder.buildContainer(); err != nil {
		return nil, err
	}

	appBuilder.ProvideModule(module.ModuleFunc(
		"application",
		ModuleFunc(appBuilder),
		provider.ProvideFunc(func() configloaderContracts.ConfigLoader { return configLoader }),
		provider.ProvideFunc(func() environment.Environment { return env }),
		provider.ProvideFunc(
			func() (loggerConfig.Registry, loggerContracts.Logger) { return loggerCfgRegistry, appLogger },
		),
		provider.ProvideFunc(func() appConfig.Config { return appCfg }),
	))

	return appBuilder, nil
}

func (b *applicationBuilder) ProvideModule(modules ...module.Module) {
	b.modules = append(b.modules, modules...)
}

func (b *applicationBuilder) Build() contracts.Application {
	return NewApplication(
		b.modules,
		b.providers,
		b.decorators,
		b.container,
		b.logger,
		b.environment,
	)
}

func (b *applicationBuilder) buildContainer() error {
	cfg, err := containerConfig.Load(b.configLoader, b.environment)
	if err != nil {
		return syserrors.Wrap(err, "failed to load container config")
	}

	log := b.logger
	if cfg.Logger != "" {
		loggerCfg, err := b.loggerConfigRegistry.GetByName(cfg.Logger)
		if err != nil {
			return syserrors.Wrapf(err, "failed to load container logger config %s", cfg.Logger)
		}

		log, err = logger.GetInstance(&loggerCfg)
		if err != nil {
			return syserrors.Wrapf(err, "failed to initialize container logger %s", cfg.Logger)
		}
	}

	w, err := waiter.New(log, waiterConfig.CancelOnShutdownSignal(true))
	if err != nil {
		return err
	}

	adapter, err := dig.New(&cfg, log, w)
	if err != nil {
		return err
	}

	b.container = container.New(adapter)

	return nil
}

func (b *applicationBuilder) Modules() []module.Module                         { return b.modules }
func (b *applicationBuilder) Providers() []provider.Provider                   { return b.providers }
func (b *applicationBuilder) Decorators() []decorator.Decorator                { return b.decorators }
func (b *applicationBuilder) Container() containerContracts.Container          { return b.container }
func (b *applicationBuilder) Logger() loggerContracts.Logger                   { return b.logger }
func (b *applicationBuilder) ConfigLoader() configloaderContracts.ConfigLoader { return b.configLoader }
func (b *applicationBuilder) Environment() environment.Environment             { return b.environment }
