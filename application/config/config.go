package config

import (
	"strings"

	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Name        string `mapstructure:"name"`
	Description string `mapstructure:"description"`
	Logger      string `mapstructure:"logger"`
}

type Option = options.OptionFn[Config]

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required),
	)
}

func Default() Config {
	return Config{}
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (Config, error) {
	var zero Config

	cfg := Default()

	if err := loader.LoadKey(env, &cfg, "app"); err != nil {
		return zero, syserrors.Wrap(err, "failed to load app config key")
	}
	if err := options.Apply(&cfg, opts...); err != nil {
		return zero, syserrors.Wrap(err, "failed to apply options on app config")
	}

	return cfg, nil
}

var Name = options.Option(func(c *Config, name string) {
	c.Name = strings.TrimSpace(name)
})

var Description = options.Option(func(c *Config, description string) {
	c.Description = strings.TrimSpace(description)
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
