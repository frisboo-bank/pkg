package config

import (
	digConfig "frisboo-bank/pkg/container/adapters/dig/config"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Type = options.Option(func(c *Config, sType containertype.ContainerType) {
	c.Type = sType
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Tracing = options.Option(func(c *Config, tracing bool) {
	c.Tracing = tracing
})

var Dig = options.Option(func(c *Config, dig *digConfig.Config) {
	c.Dig = dig
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
