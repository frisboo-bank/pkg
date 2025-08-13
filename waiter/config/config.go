package config

import (
	"context"
	"errors"

	"frisboo-bank/pkg/options"
)

type Config struct {
	ParentContext          context.Context
	CancelOnShutdownSignal bool
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(&Config{})
}

func ParentContext(parentCtx context.Context) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if parentCtx == nil {
			return errors.New("parentCtx must me set")
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
