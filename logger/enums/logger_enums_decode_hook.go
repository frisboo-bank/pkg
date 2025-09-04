package enums

import (
	"fmt"
	"reflect"
	"strings"

	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"

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
		case reflect.TypeOf(loglevel.LogLevel{}):
			data, _ = strings.CutSuffix(strings.ToLower(data.(string)), "level")
			return loglevel.ParseLogLevel(fmt.Sprintf("%sLevel", data))
		case reflect.TypeOf(loggertype.LoggerType{}):
			return loggertype.ParseLoggerType(data)
		}

		return data, nil
	}
}
