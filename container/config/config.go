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
	Type    containertype.ContainerType `mapstructure:"type"`
	Debug   bool                        `mapstructure:"debug"`
	Tracing bool                        `mapstructure:"tracing"`

	// adapters
	Dig *digConfig.Config `mapstructure:"dig"`

	// dependency
	Logger *loggerConfig.Config `mapstructure:"logger"`
}

func Default() *Config {
	return &Config{
		Type:    containertype.ContainerTypes.DIG,
		Debug:   false,
		Tracing: false,
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
