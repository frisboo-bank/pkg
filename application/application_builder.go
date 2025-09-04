package application

import (
	"fmt"
	"os"

	appConfig "frisboo-bank/pkg/application/config"
	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/application/infrastructure"
	configloader "frisboo-bank/pkg/config/config_loader"
	configloaderConfig "frisboo-bank/pkg/config/config_loader/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/container"
	containerConfig "frisboo-bank/pkg/container/config"
	containerContracts "frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	containerEnums "frisboo-bank/pkg/container/enums"
	"frisboo-bank/pkg/environment"
	httpServerEnums "frisboo-bank/pkg/http/http_server/enums"
	"frisboo-bank/pkg/logger"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerEnums "frisboo-bank/pkg/logger/enums"
	rpcServerEnums "frisboo-bank/pkg/rpc/rpc_server/enums"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ contracts.ApplicationBuilder = (*applicationBuilder)(nil)

type applicationBuilder struct {
	container   containerContracts.Container
	providers   []provider.Provider
	decorators  []decorator.Decorator
	modules     []module.Module
	logger      loggerContracts.Logger
	environment environment.Environment
}

func NewApplicationBuilder(environments ...environment.Environment) contracts.ApplicationBuilder {
	env := environment.Load(environments...)

	configLoader, err := getConfigLoader()
	if err != nil {
		fmt.Println(
			syserrors.Message(
				syserrors.Newf("failed to instantiate the config loader: got %w", err),
				[]string{"application-builder"},
			),
		)
		os.Exit(1)
	}

	appCfg, containerCfg, err := loadConfigs(configLoader, env)
	if err != nil {
		fmt.Println(
			syserrors.Message(syserrors.Newf("failed to load options: got %w", err), []string{"application-builder"}),
		)
		os.Exit(1)
	}

	appLogger, err := logger.GetInstance(&appCfg.Logger)
	if err != nil {
		fmt.Println(
			syserrors.Message(
				syserrors.Newf("failed to instantiate the app logger: got %w", err),
				[]string{"application-builder"},
			),
		)
		os.Exit(1)
	}

	containerLogger, err := logger.GetInstance(containerCfg.Logger)
	if err != nil {
		fmt.Println(
			syserrors.Message(
				syserrors.Newf("failed to instantiate the container logger: got %w", err),
				[]string{"application-builder"},
			),
		)
		os.Exit(1)
	}

	diContainer, err := container.GetInstance(containerCfg, containerLogger)
	if err != nil {
		fmt.Println(
			syserrors.Message(
				syserrors.Newf("failed to instantiate the container: got %w", err),
				[]string{"application-builder"},
			),
		)
		os.Exit(1)
	}

	return &applicationBuilder{
		modules: []module.Module{
			environment.ModuleFunc(env),
			// logger.ModuleFunc(loggerOpts),
			infrastructure.ModuleFunc(appCfg),
		},
		// providers: []provider.Provider{
		// 	provider.ProvideFunc(func() *appConfig.EnvConfig { return appOpts }),
		// 	provider.ProvideFunc(func() configloaderContracts.ConfigLoader { return configLoader }),
		// },
		decorators:  []decorator.Decorator{},
		container:   diContainer,
		logger:      appLogger,
		environment: env,
	}
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

func (b *applicationBuilder) Modules() []module.Module {
	return b.modules
}

func (b *applicationBuilder) Providers() []provider.Provider {
	return b.providers
}

func (b *applicationBuilder) Decorators() []decorator.Decorator {
	return b.decorators
}

func (b *applicationBuilder) Container() containerContracts.Container {
	return b.container
}

func (b *applicationBuilder) Logger() loggerContracts.Logger {
	return b.logger
}

func (b *applicationBuilder) Environment() environment.Environment {
	return b.environment
}

func getConfigLoader() (configloaderContracts.ConfigLoader, error) {
	cfg, err := configloaderConfig.New(
		configloaderConfig.Debug(false),

		configloaderConfig.DecodeHookFuncs(
			containerEnums.ContainerEnumsDecodeHook(),
			loggerEnums.LoggerEnumsDecodeHook(),
			httpServerEnums.HTTPServerEnumsDecodeHook(),
			rpcServerEnums.RPCServerEnumsDecodeHook(),
		),
	)
	if err != nil {
		return nil, err
	}

	return configloader.New(cfg), nil
}

func loadConfigs(
	configLoader configloaderContracts.ConfigLoader,
	env environment.Environment,
) (*appConfig.Config, *containerConfig.Config, error) {
	var errs *multierror.Error

	appOpts, err := appConfig.Load(configLoader, env)
	errs = multierror.Append(errs, syserrors.Message(err, []string{"app-config"}))

	containerOpts, err := containerConfig.Load(configLoader, env)
	errs = multierror.Append(errs, syserrors.Message(err, []string{"container-config"}))

	return appOpts, containerOpts, errs.ErrorOrNil()
}
