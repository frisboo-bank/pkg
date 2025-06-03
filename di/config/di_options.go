package config

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type DiOptions struct {
	Type   di.DiType `mapstructure:"type"`
	Logger logger.Logger
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[DiOptions]())

func ProvideLogConfig(env environment.Environment) (*DiOptions, error) {
	return config.BindConfigKey[*DiOptions](optionName, env)
}
