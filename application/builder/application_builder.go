package builder

import (
	"frisboo-bank/pkg/application/app"
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
	databaseclientEnums "frisboo-bank/pkg/database/database_client/enums"
	"frisboo-bank/pkg/environment"
	httpServerEnums "frisboo-bank/pkg/http/http_server/enums"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerEnums "frisboo-bank/pkg/logger/enums"
	migrationEnums "frisboo-bank/pkg/migration/enums"
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
			databaseclientEnums.DatabaseClientEnumsDecodeHook(),
			environment.EnvironmentEnumsDecodeHook(),
			httpServerEnums.HTTPServerEnumsDecodeHook(),
			loggerEnums.LoggerEnumsDecodeHook(),
			migrationEnums.MigrationEnumsDecodeHook(),
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
	appLogger, err := logger.GetByNameOrDefault(loggerCfgRegistry, appCfg.Logger)
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
		"core",
		provider.ProvideFunc(func() configloaderContracts.ConfigLoader { return configLoader }),
		provider.ProvideFunc(func() environment.Environment { return env }),
		provider.ProvideFunc(func() loggerConfig.Registry { return loggerCfgRegistry }),
		provider.ProvideFunc(func() loggerContracts.Logger { return appLogger }),
		provider.ProvideFunc(func() appConfig.Config { return appCfg }),
	))

	return appBuilder, nil
}

func (b *applicationBuilder) ProvideModule(modules ...module.Module) {
	b.modules = append(b.modules, modules...)
}

func (b *applicationBuilder) Build() contracts.Application {
	return app.NewApplication(
		b.container,
		b.logger,
		b.environment,
		b.modules,
		b.providers,
		b.decorators,
	)
}

func (b *applicationBuilder) buildContainer() error {
	cfg, err := containerConfig.Load(b.configLoader, b.environment)
	if err != nil {
		return syserrors.Wrap(err, "failed to load container config")
	}

	log, err := logger.GetByNameWithFallback(b.loggerConfigRegistry, cfg.Logger, b.logger)
	if err != nil {
		return syserrors.Wrapf(err, "container logger")
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
