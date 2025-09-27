package config

import (
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type  cachetype.CacheType `mapstructure:"type"`
	Debug bool                `mapstructure:"debug"`

	// adapter
	Redis redisConfig.Config `mapstructure:"redis"`

	// dependency
	Logger string `mapstructure:"logger"`
}

func Default() Config {
	return Config{
		Type:  cachetype.CacheTypes.REDIS,
		Debug: false,
		Redis: *redisConfig.Default(),
	}
}

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"caches",
		"cache",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load cache registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(cachetype.CacheTypes))),
	); err != nil {
		return err
	}

	switch c.Type {
	case cachetype.CacheTypes.REDIS:
		if err := validation.Validate(&c.Redis, validation.Required); err != nil {
			return err
		}
		return c.Redis.Validate()
	}

	return nil
}

var Type = options.Option(func(c *Config, sType cachetype.CacheType) {
	c.Type = sType
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var RedisConfig = options.Option(func(c *Config, redisConfig redisConfig.Config) {
	c.Redis = redisConfig
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
