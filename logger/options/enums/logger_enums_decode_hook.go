package enums

import (
	"reflect"

	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"

	"github.com/go-viper/mapstructure/v2"
)

func LoggerEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		switch t {
		case reflect.TypeOf(encodingtype.EncodingType{}):
			return encodingtype.ParseEncodingType(data)
		case reflect.TypeOf(logtype.LogType{}):
			return logtype.ParseLogType(data)
		case reflect.TypeOf(loglevel.LogLevel{}):
			return loglevel.ParseLogLevel(data)
		}

		return data, nil
	}
}
