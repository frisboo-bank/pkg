package invoker

import (
	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Info      any
	NamedDeps map[string]string
}

type Option = options.OptionFn[Config]

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
}

var Info = options.Option(func(c *Config, info any) {
	c.Info = info
})

func NamedDep(ref string, name string) Option {
	return func(c *Config) error {
		if c.NamedDeps == nil {
			c.NamedDeps = make(map[string]string, 1)
		}
		c.NamedDeps[ref] = name
		return nil
	}
}
