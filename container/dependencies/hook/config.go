package hook

import (
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	As         []any
	Export     bool
	Group      string
	LocationPC uintptr
	Name       string
	NamedDeps  map[string]string
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
	// validation.Field(&c.NamedDeps, validation.Map(
	// 	validation.Key("Iface", validation.Required),
	// 	validation.Key("Name", validation.Required),
	// )))
}
