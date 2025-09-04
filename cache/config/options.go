package config

import (
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType cachetype.CacheType) error {
	if err := validation.EnumOneOf("Type", sType, cachetype.CacheTypes); err != nil {
		return err
	}
	c.Type = sType
	return nil
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var RedisConfig = options.Option(func(c *Config, redisConfig redisConfig.Config) {
	c.Redis = redisConfig
})

var LoggerConfig = options.Option(func(c *Config, loggerConfig loggerConfig.Config) {
	c.Logger = loggerConfig
})
