package enums

import (
	"reflect"

	containertype "frisboo-bank/pkg/container/enums/container_type"

	"github.com/go-viper/mapstructure/v2"
)

func ContainerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeOf(containertype.ContainerType{}):
			return containertype.ParseContainerType(data)
		}

		return data, nil
	}
}
