package hook

import (
	"frisboo-bank/pkg/config"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	As         []any
	Export     bool
	Group      string
	LocationPC uintptr
	Name       string
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}
