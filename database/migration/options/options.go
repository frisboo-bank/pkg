package options

import (
	"errors"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
)

type CommandType = string

const (
	CommandTypeUp   = CommandType("up")
	CommandTypeDown = CommandType("down")
)

var ErrMigrationFailed = errors.New("migration: failed to run migration with error: %w")

type EnvConfig struct {
	Host         string
	Port         string
	User         string
	DBName       string
	SSLMode      bool
	Password     string
	MigrationDir string
}

func LoadEnvConfig(loader configContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.LoadConfig[EnvConfig](loader, env, "migration")
}

type Config struct{}

var defaultConfig = &Config{}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}

func FromEnvConfig(cfg *EnvConfig) *options.OptionBuilder[Config] {
	opts := Apply()

	return opts
}
