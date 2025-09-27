package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	postgresConfig "frisboo-bank/pkg/database/database_client/adapters/postgres/config"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
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

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"dbClients",
		"dbClient",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load database-client registry")
	}
	return reg, nil
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
