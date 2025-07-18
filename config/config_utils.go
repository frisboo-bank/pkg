package config

import (
	"frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

func LoadOptions[T any](loader contracts.ConfigLoader, env environment.Environment) (*T, error) {
	optionName := strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[T]())

	cfg := new(T)

	err := loader.LoadConfigByKey(optionName, env, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
