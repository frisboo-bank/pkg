package enums

import (
	"fmt"
	"reflect"
	"strings"

	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

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
			strData, ok := data.(string)
			if !ok {
				return nil, fmt.Errorf("expected string for loglevel, got %T", data)
			}
			strData, _ = strings.CutSuffix(strings.ToLower(strData), "level")
			return loglevel.ParseLogLevel(fmt.Sprintf("%sLevel", strData))
		case reflect.TypeOf(loggertype.LoggerType{}):
			return loggertype.ParseLoggerType(data)
		}

		return data, nil
	}
}
