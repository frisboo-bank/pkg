package config

import (
	"strings"

	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Name        string              `mapstructure:"name"`
	Description string              `mapstructure:"description"`
	Logger      loggerConfig.Config `mapstructure:"logger"`
}

func Default() *Config {
	loggerCfg := loggerConfig.Default()
	loggerCfg.Prefix = "application"

	return &Config{
		Name:        "",
		Description: "",
		Logger:      *loggerCfg,
	}
}

func (c Config) Validate() error {
	var errs *multierror.Error

	if strings.TrimSpace(c.Name) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("Name"))
	}
	errs = multierror.Append(errs, c.Logger.Validate())

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("app", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
