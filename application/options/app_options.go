package options

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type AppOptions struct {
	Name             string `mapstructure:"name"`
	EnableHttpServer bool   `mapstructure:"enableHttpServer"`
	EnableGRPCServer bool   `mapstructure:"enableGRPCServer"`
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[AppOptions]())

func ProvideLogOptions(env environment.Environment) (*AppOptions, error) {
	return config.BindConfigKey[*AppOptions](optionName, env)
}
