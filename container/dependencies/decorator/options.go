package decorator

import (
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var BeforeCallback = options.OptionErr(func(c *Config, cb CallbackFn) error {
	if cb == nil {
		return syserrors.CantBeNilError("Callback")
	}
	c.BeforeCallback = cb
	return nil
})

var Callback = options.OptionErr(func(c *Config, cb CallbackFn) error {
	if cb == nil {
		return syserrors.CantBeNilError("Callback")
	}
	c.Callback = cb
	return nil
})

var Info = options.OptionErr(func(c *Config, info any) error {
	if info == nil {
		return syserrors.CantBeNilError("Callback")
	}
	c.Info = info
	return nil
})
