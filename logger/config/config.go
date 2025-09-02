package config

import (
	"io"
	"os"

	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"

	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	CallDepth     int                       `mapstructure:"callDepth"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Output        io.Writer                 `mapstructure:"-"`
	Prefix        string                    `mapstructure:"prefix"`
	TracerEnabled bool                      `mapstructure:"tracerEnabled"`
	Type          loggertype.LoggerType     `mapstructure:"type"`
}

func Default() *Config {
	return &Config{
		CallDepth:     0,
		CallerEnabled: false,
		Encoding:      encodingtype.EncodingTypes.TEXT,
		Level:         loglevel.LogLevels.ERRORLEVEL,
		Output:        os.Stdout,
		Prefix:        "core",
		TracerEnabled: false,
		Type:          loggertype.LoggerTypes.LOGRUS,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if c.CallDepth < 0 {
		errs = multierror.Append(errs, syserrors.New("CallDepth must be a positive number: get %q", c.CallDepth))
	}
	if c.Encoding == encodingtype.EncodingTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.New("Encoding is invalid"))
	}
	if c.Level == loglevel.LogLevels.UNKNOWNLEVEL {
		errs = multierror.Append(errs, syserrors.New("Level is invalid"))
	}
	if c.Output == nil {
		errs = multierror.Append(errs, syserrors.New("Output is invalid"))
	}
	if c.Type == loggertype.LoggerTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.New("Type is invalid"))
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
