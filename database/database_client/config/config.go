package config

import (
	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/enums/database_client_type"
	"frisboo-bank/pkg/environment"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type DBClientConfig struct {
	Type     databaseclienttype.DatabaseClientType `mapstructure:"type"`
	Enabled  bool                                  `mapstructure:"enabled"`
	Host     string                                `mapstructure:"host"`
	Port     string                                `mapstructure:"port"`
	User     string                                `mapstructure:"user"`
	Password string                                `mapstructure:"password"`
	SSLMode  bool                                  `mapstructure:"sslMode"`
}

type Config struct {
	Instances []DBClientConfig `mapstructure:"instances"`

	// dependency
	Logger *loggerConfig.Config `mapstructure:"logger"`
}

func Default() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("database", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
