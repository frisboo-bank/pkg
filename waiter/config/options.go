package config

import (
	"context"
	"time"

	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var ParentContext = options.Option(func(c *Config, parentCtx context.Context) {
	c.ParentContext = parentCtx
})

var CancelOnShutdownSignal = options.Option(func(c *Config, cancelOnShutdownSignal bool) {
	c.CancelOnShutdownSignal = cancelOnShutdownSignal
})

// WaitTimeout sets the maximum duration for each Wait hook.
var WaitTimeout = options.Option(func(c *Config, waitTimeout time.Duration) {
	c.WaitTimeout = waitTimeout
})

// CleanupTimeout sets the maximum duration for each Cleanup hook.
var CleanupTimeout = options.Option(func(c *Config, cleanupTimeout time.Duration) {
	c.CleanupTimeout = cleanupTimeout
})
