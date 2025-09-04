package config

import (
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Type  cachetype.CacheType `mapstructure:"type"`
	Debug bool                `mapstructure:"debug"`

	// adapter
	Redis redisConfig.Config `mapstructure:"redis"`

	// dependency
	Logger loggerConfig.Config `mapstructure:"logger"`
}

func Default() *Config {
	return &Config{
		Debug: false,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	switch c.Type {
	case cachetype.CacheTypes.REDIS:
		errs = multierror.Append(errs, c.Redis.Validate())
	}

	errs = multierror.Append(errs,
		validation.NotNil("Logger", c.Logger),
	)

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("cache", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
