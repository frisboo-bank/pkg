package config

import (
	"io"
	"os"

	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/config"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	Type          loggertype.LoggerType     `mapstructure:"type"`
	CallDepth     int                       `mapstructure:"callDepth"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Prefix        string                    `mapstructure:"prefix"`
	TracerEnabled bool                      `mapstructure:"tracerEnabled"`

	// adapter
	Logrus  *logrusConfig.Config  `mapstructure:"logrus"`
	Zerolog *zerologConfig.Config `mapstructure:"zerolog"`

	// dependency
	Output io.Writer `mapstructure:"-"`
}

func Default() *Config {
	return &Config{
		Type:          loggertype.LoggerTypes.LOGRUS,
		CallDepth:     0,
		CallerEnabled: false,
		Encoding:      encodingtype.EncodingTypes.TEXT,
		Level:         loglevel.LogLevels.ERRORLEVEL,
		Prefix:        "core",
		TracerEnabled: false,
		Output:        os.Stdout,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	errs = multierror.Append(errs,
		validation.EnumOneOf("Type", c.Type, loggertype.LoggerTypes),
		validation.EnumOneOf("Encoding", c.Encoding, encodingtype.EncodingTypes),
		validation.EnumOneOf("Level", c.Level, loglevel.LogLevels),
		validation.NonNegative("CallDepth", c.CallDepth),
		validation.NotNil("Output", c.Output),
	)

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("logger", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
