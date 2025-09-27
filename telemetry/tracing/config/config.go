package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Enabled bool `mapstructure:"enabled"`
}

func Default() Config {
	return Config{
		Enabled: true,
	}
}

type Option = options.OptionFn[Config]

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (Registry, error) {
	reg, err := registry.Load(
		configLoader,
		env,
		"telemetry.tracing",
		"telemetry.tracing",
		Default,
	)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to load telemetry tracing registry")
	}
	return reg, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}

var Enabled = options.Option(func(c *Config, enabled bool) {
	c.Enabled = enabled
})
