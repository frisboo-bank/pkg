package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/validation"

	"github.com/hashicorp/go-multierror"
)

type Config struct {
	Servers     map[string]*HTTPServerConfig `mapstructure:"servers"`
	Diagnostics bool
}

func Default() *Config {
	return &Config{
		Servers: map[string]*HTTPServerConfig{
			"default": DefaultHTTPServerConfig(),
		},
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if len(c.Servers) == 0 {
		return nil
	}

	for name, sc := range c.Servers {
		errs = multierror.Append(errs,
			validation.NotNil(name, sc),
			sc.Validate(),
		)
	}

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("http-server", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
