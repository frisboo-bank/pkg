package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
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

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"telemetry.tracing",
		"telemetry.tracing",
		Default,
	)
}
