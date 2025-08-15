package contracts

import (
	"context"

	"frisboo-bank/pkg/container/dependencies"
)

type Application interface {
	ResolveFunc(invoker dependencies.Invoker)
	RegisterHook(hook dependencies.Hooks)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
