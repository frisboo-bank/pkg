package config

import (
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
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

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(c.Type, validation.Required, validation.By(cValidation.EnumOneOf(cachetype.CacheTypes))),
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

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"caches",
		"cache",
		Default,
	)
}
