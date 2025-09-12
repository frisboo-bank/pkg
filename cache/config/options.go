package config

import (
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

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
