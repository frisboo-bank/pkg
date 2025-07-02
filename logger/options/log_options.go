package options

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type (
	LogType  string
	LogLevel string
)

const (
	TypeLogrus = LogType("logrus")
	TypeNoop   = LogType("noop")
)

const (
	LogLevelDebug = LogLevel("debug")
	LogLevelInfo  = LogLevel("info")
	LogLevelWarn  = LogLevel("warn")
	LogLevelError = LogLevel("error")
	LogLevelPanic = LogLevel("panic")
	LogLevelFatal = LogLevel("fatal")
)

type LogOptions struct {
	Level         LogLevel `mapstructure:"level"`
	Type          LogType  `mapstructure:"type"`
	CallerEnabled bool     `mapstructure:"callerEnabled"`
	EnableTracing bool     `mapstructure:"enableTracing"`
}

type LogOption func(config *LogOptions)

var DefaultLogOptions = LogOptions{
	Level: LogLevelError,
	Type:  TypeLogrus,
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[LogOptions]())

func ProvideLogOptions(env environment.Environment) (*LogOptions, error) {
	opt, err := config.BindConfigKey[*LogOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	return setDefaultOptions(opt), nil
}

func setDefaultOptions(options *LogOptions) *LogOptions {
	if options.Level == "" {
		options.Level = DefaultLogOptions.Level
	}

	if options.Type == "" {
		options.Type = DefaultLogOptions.Type
	}

	return options
}
