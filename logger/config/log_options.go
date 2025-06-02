package config

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type LogOptions struct {
	Level         logger.LogLevel `mapstructure:"level"         default:"logrus"`
	Type          logger.LogType  `mapstructure:"type"`
	CallerEnabled bool            `mapstructure:"callerEnabled"`
	EnableTracing bool            `mapstructure:"enableTracing" default:"true"`
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[LogOptions]())

func ProvideLogConfig(env environment.Environment) (*LogOptions, error) {
	return config.BindConfigKey[*LogOptions](optionName, env)
}
