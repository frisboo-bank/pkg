package config

import (
	"frisboo-bank/pkg/config"
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	digConfig "frisboo-bank/pkg/container/adapters/dig/config"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	"frisboo-bank/pkg/environment"
	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Debug   bool                        `mapstructure:"debug"`
	Logger  loggerConfig.Config         `mapstructure:"logger"`
	Tracing bool                        `mapstructure:"tracing"`
	Type    containertype.ContainerType `mapstructure:"type"`
	Dig     digConfig.Config            `mapstructure:"dig"`
}

func Default() *Config {
	loggerCfg := loggerConfig.Default()
	loggerCfg.Prefix = "container"

	return &Config{
		Debug:   false,
		Dig:     *digConfig.Default(),
		Logger:  *loggerCfg,
		Tracing: false,
		Type:    containertype.ContainerTypes.DIG,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	errs = multierror.Append(errs, c.Logger.Validate())

	if c.Type == containertype.ContainerTypes.UNKNOWN {
		errs = multierror.Append(errs, syserrors.UnknownEnumError("Type", containertype.ContainerTypes.All()))
	}

	switch c.Type {
	case containertype.ContainerTypes.DIG:
		errs = multierror.Append(errs, c.Dig.Validate())
	}

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("container", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
