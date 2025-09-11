package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	DBKey         string `mapstructure:"db"`
	MigrationsDir string `mapstructure:"migrationsDir"`
}

func Default() Config {
	return Config{
		MigrationsDir: "./migrations",
	}
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(c.DBKey, validation.Required),
		validation.Field(c.MigrationsDir, validation.Required),
	)
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"migrations",
		"migrations",
		Default,
	)
}
