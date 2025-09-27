package config

import (
	"io"
	"os"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/config"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type          loggertype.LoggerType     `mapstructure:"type"`
	CallDepth     int                       `mapstructure:"callDepth"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Prefix        string                    `mapstructure:"prefix"`
	TracerEnabled bool                      `mapstructure:"tracerEnabled"`

	// adapters
	Logrus  logrusConfig.Config  `mapstructure:"logrus"`
	Zerolog zerologConfig.Config `mapstructure:"zerolog"`

	// dependencies
	Output io.Writer `mapstructure:"-"`
}

func Default() Config {
	return Config{
		Type:          loggertype.LoggerTypes.LOGRUS,
		CallDepth:     0,
		CallerEnabled: false,
		Encoding:      encodingtype.EncodingTypes.TEXT,
		Level:         loglevel.LogLevels.ERRORLEVEL,
		Prefix:        "core",
		TracerEnabled: false,

		Logrus:  logrusConfig.Default(),
		Zerolog: zerologConfig.Default(),

		Output: os.Stdout,
	}
}

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"loggers",
		"logger",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load logger registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(loggertype.LoggerTypes))),
		validation.Field(&c.CallDepth, validation.Min(0)),
		validation.Field(&c.Encoding, validation.Required, validation.By(cValidation.EnumOneOf(encodingtype.EncodingTypes))),
		validation.Field(&c.Level, validation.Required, validation.By(cValidation.EnumOneOf(loglevel.LogLevels))),
		validation.Field(&c.Prefix, validation.Required),
		validation.Field(&c.Output, validation.Required),
	); err != nil {
		return err
	}

	switch c.Type {
	case loggertype.LoggerTypes.LOGRUS:
		if err := validation.Validate(&c.Logrus, validation.Required); err != nil {
			return err
		}
		return c.Logrus.Validate()
	case loggertype.LoggerTypes.ZEROLOG:
		if err := validation.Validate(&c.Zerolog, validation.Required); err != nil {
			return err
		}
		return c.Zerolog.Validate()
	}

	return nil
}

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

var Logrus = options.Option(func(c *Config, logrus logrusConfig.Config) {
	c.Logrus = logrus
})

var Zerolog = options.Option(func(c *Config, zerolog zerologConfig.Config) {
	c.Zerolog = zerolog
})

var Output = options.Option(func(c *Config, output io.Writer) {
	c.Output = output
})
