package config

import (
	inMemoryConfig "frisboo-bank/pkg/cache/adapters/in_memory/config"
	redisConfig "frisboo-bank/pkg/cache/adapters/redis/config"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var InMemoryConfig = options.Option(func(c *Config, inMemoryConfig inMemoryConfig.Config) {
	c.InMemory = inMemoryConfig
})

var LoggerConfig = options.Option(func(c *Config, loggerConfig loggerConfig.Config) {
	c.Logger = loggerConfig
})

var RedisConfig = options.Option(func(c *Config, redisConfig redisConfig.Config) {
	c.Redis = redisConfig
})

var Type = options.OptionErr(func(c *Config, sType cachetype.CacheType) error {
	if sType == cachetype.CacheTypes.UNKNOWN {
		return syserrors.UnknownEnumError("Type", cachetype.CacheTypes.All())
	}
	c.Type = sType
	return nil
})
