package config

import (
	"io"
	"os"

	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	logrusConfig "frisboo-bank/pkg/logger/adapters/logrus/Config"
	zerologConfig "frisboo-bank/pkg/logger/adapters/zerolog/Config"
	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

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

	// adapters
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

	if c.Type == loggertype.LoggerTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("Type", loggertype.LoggerTypes.All()))
	}
	if c.CallDepth < 0 {
		errs = multierror.Append(errs, syserrors.CantBeNegativeError("CallDepth", c.CallDepth))
	}
	if c.Encoding == encodingtype.EncodingTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("Encoding", encodingtype.EncodingTypes.All()))
	}
	if c.Level == loglevel.LogLevels.UNKNOWNLEVEL {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("Level", loglevel.LogLevels.All()))
	}
	if c.Output == nil {
		errs = multierror.Append(errs, syserrors.CantBeNilError("Output"))
	}

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
