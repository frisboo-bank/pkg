package config

import (
	"time"

	cleanuppolicy "frisboo-bank/pkg/cache/enums/cleanup_policy"
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	AllowStaleOnExpire     bool                        `mapstructure:"allowStaleOnExpire"`
	CleanupInterval        time.Duration               `mapstructure:"cleanupInterval"`
	CleanupPolicy          cleanuppolicy.CleanupPolicy `mapstructure:"cleanuppolicy"`
	DefaultTTL             time.Duration               `mapstructure:"defaultTTL"`
	DisableCleanup         bool                        `mapstructure:"disableCleanup"`
	MaxEntries             int                         `mapstructure:"maxEntries"`
	MaxMemoryBytes         int64                       `mapstructure:"maxMemoryBytes"`
	MaxMemoryBytesPerValue int64                       `mapstructure:"MaxMemoryBytesPerValue"`
	SoftTTL                time.Duration               `mapstructure:"softTTL"`
}

func Default() *Config {
	return &Config{
		AllowStaleOnExpire:     false,
		CleanupInterval:        1 * time.Minute,
		CleanupPolicy:          cleanuppolicy.CleanupPolicies.LRU,
		DefaultTTL:             5 * time.Minute,
		DisableCleanup:         false,
		MaxEntries:             100_000,
		MaxMemoryBytes:         0,
		MaxMemoryBytesPerValue: 0,
		SoftTTL:                0,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if c.CleanupInterval < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("CleanupInterval", c.CleanupInterval))
	}
	if c.CleanupPolicy == cleanuppolicy.CleanupPolicies.UNKNOWN {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("CleanupPolicy", cleanuppolicy.CleanupPolicies.All()))
	}
	if c.DefaultTTL < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("DefaultTTL", c.DefaultTTL))
	}
	if c.MaxEntries < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("MaxEntries", c.MaxEntries))
	}
	if c.MaxMemoryBytes < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("MaxMemoryBytes", c.MaxMemoryBytes))
	}
	if c.MaxMemoryBytesPerValue < 0 {
		errs = multierror.Append(
			errs,
			syserrors.CantBeNegativeError("MaxMemoryBytesPerValue", c.MaxMemoryBytesPerValue),
		)
	}
	if c.SoftTTL < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("SoftTTL", c.SoftTTL))
	}

	return errs.ErrorOrNil()
}
