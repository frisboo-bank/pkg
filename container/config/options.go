package config

import (
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Logger = options.Option(func(c *Config, logger loggerConfig.Config) {
	c.Logger = logger
})

var Tracing = options.Option(func(c *Config, tracing bool) {
	c.Tracing = tracing
})
