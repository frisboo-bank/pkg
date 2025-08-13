package config

import (
	"io"
	"os"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"

	configContracts "frisboo-bank/pkg/config/contracts"

	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

type EnvConfig struct {
	CallDepth     int                       `mapstructure:"callDepth"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	TracerEnabled bool                      `mapstructure:"tracerEnabled"`
	Type          loggertype.LoggerType     `mapstructure:"type"`
}

func LoadEnvConfig(loader configContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.LoadConfig[EnvConfig](loader, env, "logger")
}

type Config struct {
	CallDepth     int
	CallerEnabled bool
	Encoding      encodingtype.EncodingType
	Level         loglevel.LogLevel
	Name          string
	Output        io.Writer
	Prefix        string
	TracerEnabled bool
}

var defaultConfig = &Config{
	CallDepth:     4,
	CallerEnabled: true,
	Encoding:      encodingtype.EncodingTypes.TEXT,
	Level:         loglevel.LogLevels.ERRORLEVEL,
	Output:        os.Stdout,
	TracerEnabled: true,
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}

func FromEnvConfig(cfg *EnvConfig) *options.OptionBuilder[Config] {
	opts := Apply()

	if cfg.CallDepth != 0 {
		opts.With(CallDepth(cfg.CallDepth))
	}

	opts.With(CallerEnabled(cfg.CallerEnabled))

	if cfg.Encoding != encodingtype.EncodingTypes.UNKNOWN {
		opts.With(Encoding(cfg.Encoding))
	}

	if cfg.Level != loglevel.LogLevels.UNKNOWNLEVEL {
		opts.With(Level(cfg.Level))
	}

	opts.With(TracerEnabled(cfg.TracerEnabled))

	return opts
}

func CallDepth(callDepth int) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.CallDepth = callDepth
		return nil
	})
}

func CallerEnabled(CallerEnabled bool) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.CallerEnabled = CallerEnabled
		return nil
	})
}

func Encoding(encoding encodingtype.EncodingType) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.Encoding = encoding
		return nil
	})
}

func Level(level loglevel.LogLevel) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.Level = level
		return nil
	})
}

func Output(output io.Writer) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.Output = output
		return nil
	})
}

func Prefix(prefix string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.Prefix = prefix
		return nil
	})
}

func TracerEnabled(TracerEnabled bool) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		cfg.TracerEnabled = TracerEnabled
		return nil
	})
}
