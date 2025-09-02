package config

import (
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	"io"

	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/Config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/Config"
	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType loggertype.LoggerType) error {
	if sType == loggertype.LoggerTypes.UNKNOWN {
		return syserrors.UnknownEnumError("Type", loggertype.LoggerTypes.All())
	}
	c.Type = sType
	return nil
})

var CallDepth = options.OptionErr(func(c *Config, callDepth int) error {
	if callDepth < 0 {
		return syserrors.CantBeNegativeError("CallDepth", callDepth)
	}
	c.CallDepth = callDepth
	return nil
})

var CallerEnabled = options.Option(func(c *Config, CallerEnabled bool) {
	c.CallerEnabled = CallerEnabled
})

var Encoding = options.OptionErr(func(c *Config, encoding encodingtype.EncodingType) error {
	if encoding == encodingtype.EncodingTypes.UNKNOWN {
		return syserrors.UnknownEnumError("Encoding", encodingtype.EncodingTypes.All())
	}
	c.Encoding = encoding
	return nil
})

var Level = options.OptionErr(func(c *Config, level loglevel.LogLevel) error {
	if level == loglevel.LogLevels.UNKNOWNLEVEL {
		return syserrors.UnknownEnumError("Level", loglevel.LogLevels.All())
	}
	c.Level = level
	return nil
})

var Prefix = options.Option(func(c *Config, prefix string) {
	c.Prefix = prefix
})

var TracerEnabled = options.Option(func(c *Config, TracerEnabled bool) {
	c.TracerEnabled = TracerEnabled
})

var Logrus = options.OptionErr(func(c *Config, logrus *logrusConfig.Config) error {
	if logrus == nil {
		return syserrors.CantBeNilError("Logrus")
	}
	c.Logrus = logrus
	return nil
})

var Zerolog = options.OptionErr(func(c *Config, zerolog *zerologConfig.Config) error {
	if zerolog == nil {
		return syserrors.CantBeNilError("Zerolog")
	}
	c.Zerolog = zerolog
	return nil
})

var Output = options.OptionErr(func(c *Config, output io.Writer) error {
	if output == nil {
		return syserrors.CantBeNilError("Output")
	}
	c.Output = output
	return nil
})
