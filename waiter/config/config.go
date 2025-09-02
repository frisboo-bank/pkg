package config

import (
	"context"
	"frisboo-bank/pkg/options"
	"time"
)

type Config struct {
	ParentContext          context.Context
	CancelOnShutdownSignal bool
	WaitTimeout            time.Duration
	CleanupTimeout         time.Duration
}

func Default() *Config {
	return &Config{
		ParentContext:          context.Background(),
		CancelOnShutdownSignal: false,
		WaitTimeout:            30 * time.Second,
		CleanupTimeout:         5 * time.Second,
	}
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}
