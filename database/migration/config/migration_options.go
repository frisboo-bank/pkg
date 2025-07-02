package config

import (
	"errors"

	configUtils "frisboo-bank/pkg/config/utils"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type CommandType = string

const (
	CommandTypeUp   = CommandType("up")
	CommandTypeDown = CommandType("down")
)

var ErrMigrationFailed = errors.New("migration: failed to run migration with error: %w")

type MigrationOptions struct {
	Host         string
	Port         string
	User         string
	DBName       string
	SSLMode      bool
	Password     string
	MigrationDir string
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[MigrationOptions]())

func ProvideMigrationConfig(environment environment.Environment) (*MigrationOptions, error) {
	return configUtils.BindConfigKey[*MigrationOptions](optionName, environment)
}
