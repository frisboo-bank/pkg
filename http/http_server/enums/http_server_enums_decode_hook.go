package enums

import (
	"reflect"

	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"

	"github.com/go-viper/mapstructure/v2"
)

func HTTPServerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeOf(httpservertype.HttpServerType{}):
			return httpservertype.ParseHttpServerType(data)
		}

		return data, nil
	}
}
