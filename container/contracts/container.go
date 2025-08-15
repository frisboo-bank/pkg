package contracts

import (
	"context"
	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"

	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
)

type Container interface {
	RegisterModule(modules ...module.Module) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}

type ContainerAdapter interface {
	RegisterDecorator(decorators ...decorator.Decorator) error
	RegisterHook(hooks ...hook.Hooks) error
	RegisterInvoker(invokers ...invoker.Invoker) error
	RegisterProvider(providers ...provider.Provider) error
	Setup(cfg *config.Config) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}
