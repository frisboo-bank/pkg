package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	postgresConfig "frisboo-bank/pkg/database/database_client/adapters/postgres/config"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	"frisboo-bank/pkg/environment"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationIs "github.com/go-ozzo/ozzo-validation/v4/is"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type     databaseclienttype.DatabaseClientType `mapstructure:"type"`
	Enabled  bool                                  `mapstructure:"enabled"`
	Host     string                                `mapstructure:"host"`
	Port     string                                `mapstructure:"port"`
	User     string                                `mapstructure:"user"`
	Password string                                `mapstructure:"password"`
	SSLMode  bool                                  `mapstructure:"sslMode"`

	// adapters
	Postgres postgresConfig.Config `mapstructure:"postgres"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func Default() Config {
	return Config{}
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(databaseclienttype.DatabaseClientTypes))),
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.User, validation.Required),
		validation.Field(&c.Password, validation.Required),
	); err != nil {
		return err
	}

	switch c.Type {
	case databaseclienttype.DatabaseClientTypes.POSTGRES:
		if err := validation.Validate(&c.Postgres, validation.Required); err != nil {
			return err
		}
		return c.Postgres.Validate()
	}

	return nil
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"dbClients",
		"dbClient",
		Default,
	)
}
