package application

import (
	"fmt"
	"os"

	"frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/application/infrastructure"
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/container"
	containerConfig "frisboo-bank/pkg/container/config"
	containerContracts "frisboo-bank/pkg/container/contracts"
	containerEnums "frisboo-bank/pkg/container/contracts/enums"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/environment"
	httpServerEnums "frisboo-bank/pkg/http/http_server/contracts/enums"
	"frisboo-bank/pkg/logger"
	loggerConfig "frisboo-bank/pkg/logger/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	loggerEnums "frisboo-bank/pkg/logger/contracts/enums"
	rpcServerEnums "frisboo-bank/pkg/rpc/rpc_server/contracts/enums"
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
	env := environment.GetEnvFromConfig(environments...)

	configLoader := config.NewConfigLoader().
		WithDecodeHooks(
			containerEnums.ContainerEnumsDecodeHook(),
			rpcServerEnums.RPCServerEnumsDecodeHook(),
			loggerEnums.LoggerEnumsDecodeHook(),
			httpServerEnums.HTTPServerEnumsDecodeHook(),
		)

	loggerEnvCfg, err := loggerConfig.LoadEnvConfig(configLoader, env)
	if err != nil {
		fmt.Printf("application-builder: failed to load Logger options with error: %v\n", err)
		os.Exit(1)
	}

	logger, err := logger.GetInstance(loggerEnvCfg.Type, loggerConfig.FromEnvConfig(loggerEnvCfg))
	if err != nil {
		fmt.Printf("application-builder: failed to create Logger with error: %v\n", err)
		os.Exit(1)
	}

	diContainer, err := container.GetInstance(
		containertype.ContainerTypes.DIG,
		logger,
		containerConfig.Apply(),
	)
	if err != nil {
		return nil
	}

	return &applicationBuilder{
		modules: []module.Module{
			infrastructure.Module,
		},
		providers: []provider.Provider{
			provider.Provide(func() environment.Environment { return env }),
			provider.Provide(func() configContracts.ConfigLoader { return configLoader }),
			provider.Provide(func() *loggerConfig.EnvConfig { return loggerEnvCfg }),
			provider.Provide(func() loggerContracts.Logger { return logger }),
		},
		decorators:  []decorator.Decorator{},
		container:   diContainer,
		logger:      logger,
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
