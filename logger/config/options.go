package config

import (
	"io"

	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/config"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Type = options.Option(func(c *Config, sType loggertype.LoggerType) {
	c.Type = sType
})

var CallDepth = options.Option(func(c *Config, callDepth int) {
	c.CallDepth = callDepth
})

var CallerEnabled = options.Option(func(c *Config, CallerEnabled bool) {
	c.CallerEnabled = CallerEnabled
})

var Encoding = options.Option(func(c *Config, encoding encodingtype.EncodingType) {
	c.Encoding = encoding
})

var Level = options.Option(func(c *Config, level loglevel.LogLevel) {
	c.Level = level
})

var Prefix = options.Option(func(c *Config, prefix string) {
	c.Prefix = prefix
})

var TracerEnabled = options.Option(func(c *Config, TracerEnabled bool) {
	c.TracerEnabled = TracerEnabled
})

var Zerolog = options.Option(func(c *Config, zerolog *zerologConfig.Config) {
	c.Zerolog = zerolog
})

var Output = options.Option(func(c *Config, output io.Writer) {
	c.Output = output
})
