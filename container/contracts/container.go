package contracts

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	containertype "frisboo-bank/pkg/container/enums/container_type"
)

type Container interface {
	RegisterModule(modules ...module.Module) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}

type ContainerAdapter interface {
	RegisterDecorator(name string, decorator decorator.Decorator) error
	RegisterHook(hook hook.Hooks) error
	RegisterInvoker(name string, invoker invoker.Invoker) error
	RegisterProvider(name string, provider provider.Provider) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}
