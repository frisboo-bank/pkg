package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	WaitFunc    func(ctx context.Context) error
	CleanupFunc func(ctx context.Context) error
)

type WaiterHook struct {
	Wait    WaitFunc
	Cleanup CleanupFunc
}

type Waiter interface {
	Add(hooks ...WaiterHook)
	Wait() error
	Cancel()
	Context() context.Context
	Logger() loggerContracts.Logger
}
