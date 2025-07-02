package contracts

import (
	"context"

	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type Application interface {
	ResolveFunc(invoker container.Invoker)
	RegisterHook(hook container.HookStarter)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	Logger() loggerContracts.Logger
	Environment() environment.Environment
}
