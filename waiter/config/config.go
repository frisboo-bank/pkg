package config

import (
	"context"
	"errors"
	"time"

	"frisboo-bank/pkg/options"
)

type Config struct {
	ParentContext          context.Context
	CancelOnShutdownSignal bool
	WaitTimeout            time.Duration
	CleanupTimeout         time.Duration
}

var defaultConfig = &Config{
	ParentContext:          context.Background(),
	CancelOnShutdownSignal: false,
	WaitTimeout:            30 * time.Second,
	CleanupTimeout:         5 * time.Second,
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}

func ParentContext(parentCtx context.Context) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if parentCtx == nil {
			return errors.New("parentCtx must be set")
		}

		cfg.ParentContext = parentCtx
		return nil
	})
}

func CancelOnShutdownSignal(cancelOnShutdownSignal bool) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.CancelOnShutdownSignal = cancelOnShutdownSignal
		return nil
	})
}

// WaitTimeout sets the maximum duration for each Wait hook.
// Zero means "use default"; negative means "no timeout".
func WaitTimeout(d time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.WaitTimeout = d
		return nil
	})
}

// CleanupTimeout sets the maximum duration for each Cleanup hook.
// Zero means "use default"; negative means "no timeout".
func CleanupTimeout(d time.Duration) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.CleanupTimeout = d
		return nil
	})
}
