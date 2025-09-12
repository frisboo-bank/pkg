package decorator

import (
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type CallbackFn func(args ...any) error

type Config struct {
	BeforeCallback CallbackFn
	Callback       CallbackFn
	Info           any
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}
