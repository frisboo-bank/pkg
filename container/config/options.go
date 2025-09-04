package config

import (
	digConfig "frisboo-bank/pkg/container/adapters/dig/config"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType containertype.ContainerType) error {
	if err := validation.EnumOneOf("Type", sType, containertype.ContainerTypes); err != nil {
		return err
	}
	c.Type = sType
	return nil
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Tracing = options.Option(func(c *Config, tracing bool) {
	c.Tracing = tracing
})

var Dig = options.OptionErr(func(c *Config, dig *digConfig.Config) error {
	if err := validation.NotNil("Dig", dig); err != nil {
		return err
	}
	c.Dig = dig
	return nil
})

var Logger = options.OptionErr(func(c *Config, logger *loggerConfig.Config) error {
	if err := validation.NotNil("Logger", logger); err != nil {
		return err
	}
	c.Logger = logger
	return nil
})
