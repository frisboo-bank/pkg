package contracts

import (
	"context"
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
}
