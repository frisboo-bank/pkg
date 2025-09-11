package application

import (
	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/syserrors"

	configloader "frisboo-bank/pkg/config/config_loader"
	configloaderConfig "frisboo-bank/pkg/config/config_loader/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	containerContracts "frisboo-bank/pkg/container/contracts"

	containerConfig "frisboo-bank/pkg/container/config"
	containerEnums "frisboo-bank/pkg/container/enums"

	httpServerEnums "frisboo-bank/pkg/http/http_server/enums"

	appConfig "frisboo-bank/pkg/application/config"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerEnums "frisboo-bank/pkg/logger/enums"
	rpcServerEnums "frisboo-bank/pkg/rpc/rpc_server/enums"
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
	env := environment.Load(environments...)

	configLoader, err := configloader.New(
		configloaderConfig.Debug(false),
		configloaderConfig.DecodeHookFuncs(
			containerEnums.ContainerEnumsDecodeHook(),
			loggerEnums.LoggerEnumsDecodeHook(),
			httpServerEnums.HTTPServerEnumsDecodeHook(),
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

	appModule := module.ModuleFunc(
		"app",
		provider.ProvideFunc(func() environment.Environment { return env }),
		provider.ProvideFunc(func() loggerContracts.Logger { return appLogger }),
		provider.ProvideFunc(func() configloaderContracts.ConfigLoader { return configLoader }),
		provider.ProvideFunc(func() loggerConfig.Registry { return loggerCfgRegistry }),
		provider.ProvideFunc(func() appConfig.Config { return appCfg }),
	)

	appBuilder := &applicationBuilder{
		environment: env,
		logger:      appLogger,

		configLoader:         configLoader,
		loggerConfigRegistry: loggerCfgRegistry,
		appConfig:            appCfg,

		modules:    []module.Module{appModule},
		providers:  []provider.Provider{},
		decorators: []decorator.Decorator{},
	}

	diContainer, err := appBuilder.buildContainer()
	if err != nil {
		return nil, err
	}
	appBuilder.container = diContainer

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

func (b *applicationBuilder) buildContainer() (containerContracts.Container, error) {
	configRegistry, err := containerConfig.LoadRegistry(b.configLoader, b.environment)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load container config registry")
	}

	cfg, err := configRegistry.GetByNameOrDefault(b.appConfig.Container)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load container config")
	}

	diLogger := b.logger
	if cfg.Logger != "" {
		loggerCfg, err := b.loggerConfigRegistry.GetByName(cfg.Logger)
		if err != nil {
			return nil, syserrors.Wrapf(err, "failed to load container logger config %s", cfg.Logger)
		}

		diLogger, err = logger.GetInstance(&loggerCfg)
		if err != nil {
			return nil, syserrors.Wrapf(err, "failed to initialize container logger %s", cfg.Logger)
		}
	}

	diContainer, err := container.GetInstance(&cfg, diLogger)
	if err != nil {
		return nil, err
	}

	return diContainer, nil
}

func (b *applicationBuilder) Modules() []module.Module                { return b.modules }
func (b *applicationBuilder) Providers() []provider.Provider          { return b.providers }
func (b *applicationBuilder) Decorators() []decorator.Decorator       { return b.decorators }
func (b *applicationBuilder) Container() containerContracts.Container { return b.container }
func (b *applicationBuilder) Logger() loggerContracts.Logger          { return b.logger }
func (b *applicationBuilder) Environment() environment.Environment    { return b.environment }
