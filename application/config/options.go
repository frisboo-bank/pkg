package config

import (
	"strings"

	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var Name = options.OptionErr(func(c *Config, name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return syserrors.CantBeEmptyError("Name")
	}
	c.Name = name
	return nil
})

var Description = options.Option(func(c *Config, description string) {
	c.Description = strings.TrimSpace(description)
})

var Logger = options.Option(func(c *Config, logger loggerConfig.Config) {
	c.Logger = logger
})
