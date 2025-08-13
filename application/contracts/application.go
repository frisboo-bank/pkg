package contracts

import (
	"context"

	containerContract "frisboo-bank/pkg/container/contracts"
)

type Application interface {
	ResolveFunc(invoker containerContract.Invoker)
	RegisterHook(hook containerContract.HookStarter)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
