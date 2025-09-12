package invoker

import (
	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var Info = options.Option(func(c *Config, info any) {
	c.Info = info
})
