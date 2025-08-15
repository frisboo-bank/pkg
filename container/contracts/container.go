package contracts

import (
	"context"

	"frisboo-bank/pkg/container/config"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	"frisboo-bank/pkg/container/dependencies"
)

type Container interface {
	RegisterModule(modules ...dependencies.Module) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}

type ContainerAdapter interface {
	RegisterDecorator(decorators ...dependencies.Decorator) error
	RegisterHook(hooks ...dependencies.Hooks) error
	RegisterInvoker(invokers ...dependencies.Invoker) error
	RegisterProvider(providers ...dependencies.Provider) error
	Setup(cfg *config.Config) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}
