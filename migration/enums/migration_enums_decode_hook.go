package enums

import (
	"reflect"

	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"

	"github.com/go-viper/mapstructure/v2"
)

func MigrationEnumsDecodeHook() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		switch t {
		case reflect.TypeOf(migrationcommandtype.MigrationCommandType{}):
			return migrationcommandtype.ParseMigrationCommandType(data)
		case reflect.TypeOf(migratortype.MigratorType{}):
			return migratortype.ParseMigratorType(data)
		}

		return data, nil
	}
}
