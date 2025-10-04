package config

import (
	"strings"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type          migratortype.MigratorType `mapstructure:"type"`
	DB            string                    `mapstructure:"db"`
	MigrationsDir string                    `mapstructure:"migrationsDir"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func Default() Config {
	return Config{
		Type: migratortype.MigratorTypes.GOOSE,
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
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(migratortype.MigratorTypes))),
		validation.Field(&c.DB, validation.Required),
		validation.Field(&c.MigrationsDir, validation.Required),
	)
}

var Type = options.Option(func(c *Config, sType migratortype.MigratorType) {
	c.Type = sType
})

var DB = options.Option(func(c *Config, db string) {
	c.DB = strings.TrimSpace(db)
})

var MigrationsDir = options.Option(func(c *Config, migrationsDir string) {
	c.MigrationsDir = strings.TrimSpace(migrationsDir)
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
