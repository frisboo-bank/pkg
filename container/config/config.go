package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/config/registry"
	"frisboo-bank/pkg/environment"
	cValidation "frisboo-bank/pkg/validation"

	digConfig "frisboo-bank/pkg/container/adapters/dig/config"
	containertype "frisboo-bank/pkg/container/enums/container_type"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Type    containertype.ContainerType `mapstructure:"type"`
	Debug   bool                        `mapstructure:"debug"`
	Tracing bool                        `mapstructure:"tracing"`

	// adapter
	Dig *digConfig.Config `mapstructure:"dig"`

	// dependency
	Logger string `mapstructure:"logger"`
}

func Default() Config {
	return Config{
		Type:    containertype.ContainerTypes.DIG,
		Debug:   false,
		Tracing: false,
		Dig:     digConfig.Default(),
	}
}

func (c *Config) Validate() error {
	if err := validation.ValidateStruct(c,
		validation.Field(c.Type, validation.Required, validation.By(cValidation.EnumOneOf(containertype.ContainerTypes))),
	); err != nil {
		return err
	}

	switch c.Type {
	case containertype.ContainerTypes.DIG:
		if err := validation.Validate(&c.Dig, validation.Required); err != nil {
			return err
		}
		return c.Dig.Validate()
	}

	return nil
}

type Registry = registry.Registry[Config]

func LoadRegistry(configLoader configloaderContracts.ConfigLoader, env environment.Environment) (*Registry, error) {
	return registry.Load(
		configLoader,
		env,
		"containers",
		"container",
		Default,
	)
}
