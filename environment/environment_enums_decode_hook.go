package environment

import (
	"reflect"

	"github.com/go-viper/mapstructure/v2"
)

func EnvironmentEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeOf(Environment{}):
			return ParseEnvironment(data)
		}

		return data, nil
	}
}
