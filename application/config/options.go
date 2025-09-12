package config

import (
	"strings"

	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Name = options.Option(func(c *Config, name string) {
	c.Name = strings.TrimSpace(name)
})

var Description = options.Option(func(c *Config, description string) {
	c.Description = strings.TrimSpace(description)
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})

var Container = options.Option(func(c *Config, container string) {
	c.Container = container
})
