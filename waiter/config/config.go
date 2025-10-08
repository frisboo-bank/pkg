package config

import (
	"context"
	"time"

	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	ParentContext          context.Context
	CancelOnShutdownSignal bool
	WaitTimeout            time.Duration
	CleanupTimeout         time.Duration
}

func Default() Config {
	return Config{
		CancelOnShutdownSignal: false,
		WaitTimeout:            30 * time.Second,
		CleanupTimeout:         5 * time.Second,
	}
}

type Option = options.OptionFn[Config]

func New(opts ...Option) (Config, error) {
	var zero Config

	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return zero, err
	}

	return base, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.WaitTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.CleanupTimeout, validation.Required, validation.Min(0)),
	)
}

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
