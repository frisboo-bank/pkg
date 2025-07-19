package options

import (
	"errors"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
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

func ProvideMigrationConfig(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*MigrationOptions, error) {
	return config.LoadOptions[MigrationOptions](loader, env)
}
