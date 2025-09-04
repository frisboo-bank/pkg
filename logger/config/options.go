package config

import (
	"io"

	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/config"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType loggertype.LoggerType) error {
	if err := validation.EnumOneOf("Type", sType, loggertype.LoggerTypes); err != nil {
		return err
	}
	c.Type = sType
	return nil
})

var CallDepth = options.OptionErr(func(c *Config, callDepth int) error {
	if err := validation.NonNegative("CallDepth", callDepth); err != nil {
		return err
	}
	c.CallDepth = callDepth
	return nil
})

var CallerEnabled = options.Option(func(c *Config, CallerEnabled bool) {
	c.CallerEnabled = CallerEnabled
})

var Encoding = options.OptionErr(func(c *Config, encoding encodingtype.EncodingType) error {
	if err := validation.EnumOneOf("Encoding", encoding, encodingtype.EncodingTypes); err != nil {
		return err
	}
	c.Encoding = encoding
	return nil
})

var Level = options.OptionErr(func(c *Config, level loglevel.LogLevel) error {
	if err := validation.EnumOneOf("Level", level, loglevel.LogLevels); err != nil {
		return err
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
	if err := validation.NotNil("Logrus", logrus); err != nil {
		return err
	}
	c.Logrus = logrus
	return nil
})

var Zerolog = options.OptionErr(func(c *Config, zerolog *zerologConfig.Config) error {
	if err := validation.NotNil("Zerolog", zerolog); err != nil {
		return err
	}
	c.Zerolog = zerolog
	return nil
})

var Output = options.OptionErr(func(c *Config, output io.Writer) error {
	if err := validation.NotNil("Output", output); err != nil {
		return err
	}
	c.Output = output
	return nil
})
