package config

import (
	"frisboo-bank/pkg/options"
	"time"
)

type Option = options.OptionFn[Config]

func AllowStaleOnExpire(allowStaleOnExpire bool) Option {
	return options.Option(func(c *Config) {
		c.AllowStaleOnExpire = allowStaleOnExpire
	})
}

func CleanupInterval(cleanupInterval time.Duration) Option {
	return options.Option(func(c *Config) {
		c.CleanupInterval = cleanupInterval
	})
}

func DefaultTTL(defaultTTL time.Duration) Option {
	return options.Option(func(c *Config) {
		c.DefaultTTL = defaultTTL
	})
}

func MaxEntries(maxEntries int) Option {
	return options.Option(func(c *Config) {
		c.MaxEntries = maxEntries
	})
}
