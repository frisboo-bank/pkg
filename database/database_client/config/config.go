package config

import (
	"context"
	"net"
	"time"

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
	Enabled               bool                                  `mapstructure:"enabled"`
	Type                  databaseclienttype.DatabaseClientType `mapstructure:"type"`
	Debug                 bool                                  `mapstructure:"debug"`
	Host                  string                                `mapstructure:"host"`
	Port                  string                                `mapstructure:"port"`
	Database              string                                `mapstructure:"database"`
	User                  string                                `mapstructure:"user"`
	Password              string                                `mapstructure:"password"`
	SSLMode               bool                                  `mapstructure:"sslMode"`
	Context               context.Context                       `mapstructure:"-"`
	EnableTracing         bool                                  `mapstructure:"enableTracing"`
	ConnectionTimeout     time.Duration                         `mapstructure:"connectionTimeout"`
	MaxConnectionIdleTime time.Duration                         `mapstructure:"maxConnectionIdleTime"`
	MinPoolSize           uint64                                `mapstructure:"minPoolSize"`
	MaxPoolSize           uint64                                `mapstructure:"maxPoolSize"`

	// adapters
	MongoDB  any                   `mapstructure:"mongoDB"`
	Postgres postgresConfig.Config `mapstructure:"postgres"`

	// dependencies
	Logger string `mapstructure:"logger"`
}

func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func Default() Config {
	return Config{
		Enabled:               true,
		Debug:                 false,
		EnableTracing:         true,
		ConnectionTimeout:     60 * time.Second,
		MaxConnectionIdleTime: 3 * time.Minute,
		MinPoolSize:           20,
		MaxPoolSize:           300,
	}
}

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

type Option = options.OptionFn[Config]

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(&c.Type, validation.Required, validation.By(cValidation.EnumOneOf(databaseclienttype.DatabaseClientTypes))),
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.User, validation.Required),
		validation.Field(&c.Database, validation.Required),
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

var Type = options.Option(func(c *Config, sType databaseclienttype.DatabaseClientType) {
	c.Type = sType
})

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Host = options.Option(func(c *Config, host string) {
	c.Host = host
})

var Port = options.Option(func(c *Config, port string) {
	c.Port = port
})

var Database = options.Option(func(c *Config, database string) {
	c.Database = database
})

var User = options.Option(func(c *Config, user string) {
	c.User = user
})

var Password = options.Option(func(c *Config, password string) {
	c.Password = password
})

var SSLMode = options.Option(func(c *Config, sslMode bool) {
	c.SSLMode = sslMode
})

var Context = options.Option(func(c *Config, ctx context.Context) {
	c.Context = ctx
})

var EnableTracing = options.Option(func(c *Config, enableTracing bool) {
	c.EnableTracing = enableTracing
})

var Postgres = options.Option(func(c *Config, pgConfig postgresConfig.Config) {
	c.Postgres = pgConfig
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
