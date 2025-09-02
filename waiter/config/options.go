package config

import (
	"context"
	"time"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var ParentContext = options.OptionErr(func(c *Config, parentCtx context.Context) error {
	if parentCtx == nil {
		return syserrors.CantBeNilError("ParentContext")
	}
	c.ParentContext = parentCtx
	return nil
})

var CancelOnShutdownSignal = options.Option(func(c *Config, cancelOnShutdownSignal bool) {
	c.CancelOnShutdownSignal = cancelOnShutdownSignal
})

// WaitTimeout sets the maximum duration for each Wait hook.
var WaitTimeout = options.OptionErr(func(c *Config, d time.Duration) error {
	if d < 0 {
		return syserrors.CantBeNegativeError("WaitTimeout", d)
	}
	c.WaitTimeout = d
	return nil
})

// CleanupTimeout sets the maximum duration for each Cleanup hook.
var CleanupTimeout = options.OptionErr(func(c *Config, d time.Duration) error {
	if d < 0 {
		return syserrors.CantBeNegativeError("CleanupTimeout", d)
	}
	c.CleanupTimeout = d
	return nil
})
