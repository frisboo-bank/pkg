package enums

import (
	"reflect"

	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"

	"github.com/go-viper/mapstructure/v2"
)

func HTTPServerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeOf(httpservertype.HTTPServerType{}):
			return httpservertype.ParseHTTPServerType(data)
		}

		return data, nil
	}
}
