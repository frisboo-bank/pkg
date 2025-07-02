package options

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/noop"
)

type WaiterOptions struct {
	ParentContext          context.Context
	CancelOnShutdownSignal bool
	Logger                 loggerContracts.Logger
}

type WaiterOption func(o *WaiterOptions)

// Define the waiter context
// default: default context
func WithParentContext(ctx context.Context) WaiterOption {
	return func(o *WaiterOptions) {
		o.ParentContext = ctx
	}
}

// Cancel waiter on system signals
// default: false
func UseCancelOnShutdownSignal() WaiterOption {
	return func(o *WaiterOptions) {
		o.CancelOnShutdownSignal = true
	}
}

// Define the logger
// default: NoopLogger
func WithLogger(logger loggerContracts.Logger) WaiterOption {
	return func(o *WaiterOptions) {
		o.Logger = logger
	}
}

func GetOptionsWithDefault(config ...WaiterOption) *WaiterOptions {
	cfg := &WaiterOptions{
		ParentContext:          context.Background(),
		CancelOnShutdownSignal: false,
		Logger:                 noop.NewNoopLogger(),
	}

	for _, option := range config {
		option(cfg)
	}

	return cfg
}
