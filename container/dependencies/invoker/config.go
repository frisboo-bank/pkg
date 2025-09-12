package invoker

import (
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Info any
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}
