package config

import "frisboo-bank/pkg/options"

type Option = options.OptionFn[Config]

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})
