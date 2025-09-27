package config

import (
	"strings"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
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

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (registry.Registry[Config], error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"migrations",
		"migrations",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load migration registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.DBKey, validation.Required),
		validation.Field(&c.MigrationsDir, validation.Required),
	)
}

var DB = options.Option(func(c *Config, dbKey string) {
	c.DBKey = strings.TrimSpace(dbKey)
})

var MigrationsDir = options.Option(func(c *Config, migrationsDir string) {
	c.MigrationsDir = strings.TrimSpace(migrationsDir)
})
