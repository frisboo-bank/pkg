package enums

import (
	"reflect"

	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"

	"github.com/go-viper/mapstructure/v2"
)

func DatabaseClientEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeOf(databaseclienttype.DatabaseClientType{}):
			return databaseclienttype.ParseDatabaseClientType(data)
		}

		return data, nil
	}
}
