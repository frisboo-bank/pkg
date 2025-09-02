package config

import (
	"frisboo-bank/pkg/config"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct{}

func Default() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}
