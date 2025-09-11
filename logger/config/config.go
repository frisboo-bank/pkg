package config

import (
	"io"
	"os"

	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"

	cValidation "frisboo-bank/pkg/validation"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/config"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"

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
	Logrus  *logrusConfig.Config  `mapstructure:"logrus"`
	Zerolog *zerologConfig.Config `mapstructure:"zerolog"`

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

type Registry = *registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"loggers",
		"logger",
		Default,
	)
}
