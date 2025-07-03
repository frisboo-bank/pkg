package options

import (
	"strings"

	configloader "frisboo-bank/pkg/config/config_loader"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options/contracts"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

// ProvideOptions loads and returns a config options struct for any type T.
// defaultOpt should be a T struct with default values set.
func LoadOptions[T contracts.Options](env environment.Environment, setDefaultFunc func(options T)) (T, error) {
	optionName := strcase.LowerCamelCase(strings.TrimPrefix(typemapper.GetGenericTypeNameByT[T](), "*"))

	loader := configloader.New(configloader.ConfigBinder{})

	var cfg T
	err := loader.BindConfigKey(optionName, env, &cfg)
	if err != nil {
		return *new(T), err
	}

	if setDefaultFunc != nil {
		setDefaultFunc(cfg)
	}

	return cfg, nil
}

func CloneOptions[T any](o *T) *T {
	clone := *o
	return &clone
}
