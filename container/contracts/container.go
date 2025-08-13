package contracts

import (
	"context"

	"frisboo-bank/pkg/container/config"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
)

type Container interface {
	RegisterModule(modules ...Module) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}

type ContainerAdapter interface {
	RegisterDecorator(decorators ...Decorator) error
	RegisterHook(hooks ...HookStarter) error
	RegisterInvoker(invokers ...Invoker) error
	RegisterProvider(providers ...Provider) error
	Setup(cfg *config.Config) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Type() containertype.ContainerType
}
