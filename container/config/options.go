package config

import (
	digConfig "frisboo-bank/pkg/container/adapters/dig/config"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType containertype.ContainerType) error {
	if sType == containertype.ContainerTypes.UNKNOWN {
		return syserrors.UnknownEnumError("Type", containertype.ContainerTypes.All())
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
	if dig == nil {
		return syserrors.CantBeNilError("Dig")
	}
	c.Dig = dig
	return nil
})

var Logger = options.OptionErr(func(c *Config, logger *loggerConfig.Config) error {
	if logger == nil {
		return syserrors.CantBeNilError("Logger")
	}
	c.Logger = logger
	return nil
})
