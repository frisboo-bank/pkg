package config

import (
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct{}

func Default() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}
