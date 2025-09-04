package config

import (
	inMemoryConfig "frisboo-bank/pkg/cache/adapters/in_memory/config"
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Debug    bool                  `mapstructure:"debug"`
	InMemory inMemoryConfig.Config `mapstructure:"inMemory"`
	Logger   loggerConfig.Config   `mapstructure:"logger"`
	Redis    redisConfig.Config    `mapstructure:"redis"`
	Type     cachetype.CacheType   `mapstructure:"type"`
}

func Default() *Config {
	loggerCfg := loggerConfig.Default()
	loggerCfg.Prefix = "container"

	return &Config{
		Debug:    false,
		InMemory: *inMemoryConfig.Default(),
		Logger:   *loggerCfg,
		Redis:    *redisConfig.Default(),
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if c.Type == cachetype.CacheTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("Type", cachetype.CacheTypes.All()))
	}

	switch c.Type {
	case cachetype.CacheTypes.IN_MEMORY:
		errs = multierror.Append(errs, c.InMemory.Validate())
	case cachetype.CacheTypes.REDIS:
		errs = multierror.Append(errs, c.Redis.Validate())
	}

	errs = multierror.Append(errs, c.Logger.Validate())

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
