package contracts

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type Application interface {
	ResolveFunc(invoker invoker.Invoker)
	RegisterHook(hook hook.Hooks)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Logger() loggerContracts.Logger
}
