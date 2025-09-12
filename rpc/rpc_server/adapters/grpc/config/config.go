package config

import (
	cValidation "frisboo-bank/pkg/validation"

	"github.com/hashicorp/go-multierror"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct{}

func Default() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}
