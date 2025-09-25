package contracts

import (
	"context"
)

type (
	WaitFunc    func(ctx context.Context) error
	CleanupFunc func(ctx context.Context) error
)

type WaiterHook struct {
	Name    string
	Wait    WaitFunc
	Cleanup CleanupFunc
}

type Waiter interface {
	AddHooks(hooks ...WaiterHook) error
	AddHook(hook WaiterHook) error
	Wait() error
	Cancel()
}
