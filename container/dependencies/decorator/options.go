package decorator

import (
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var BeforeCallback = options.Option(func(c *Config, cb CallbackFn) {
	c.BeforeCallback = cb
})

var Callback = options.Option(func(c *Config, cb CallbackFn) {
	c.Callback = cb
})

var Info = options.Option(func(c *Config, info any) {
	c.Info = info
})
