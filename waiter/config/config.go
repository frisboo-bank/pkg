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
		ParentContext:          context.Background(),
		CancelOnShutdownSignal: false,
		WaitTimeout:            30 * time.Second,
		CleanupTimeout:         5 * time.Second,
	}
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ParentContext, validation.Required),
		validation.Field(&c.WaitTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.CleanupTimeout, validation.Required, validation.Min(0)),
	)
}

func New(opts ...Option) (Config, error) {
	var zero Config

	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return zero, err
	}

	return base, nil
}
