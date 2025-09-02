package invoker

import (
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var Info = options.OptionErr(func(c *Config, info any) error {
	if info == nil {
		return syserrors.CantBeNilError("Info")
	}
	c.Info = info
	return nil
})
