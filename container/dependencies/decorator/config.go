package decorator

import (
	"frisboo-bank/pkg/config"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type CallbackFn func(args ...any) error

type Config struct {
	BeforeCallback CallbackFn
	Callback       CallbackFn
	Info           any
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}
