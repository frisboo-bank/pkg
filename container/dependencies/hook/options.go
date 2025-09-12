package hook

import (
	"reflect"
	"strings"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option = options.OptionFn[Config]

var As = options.VarOptionErr(func(c *Config, ifaces ...any) error {
	if len(ifaces) == 0 {
		return syserrors.CantBeEmptyError("As")
	}
	for _, i := range ifaces {
		t := reflect.TypeOf(i)
		if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Interface {
			return syserrors.New("As expects pointers to interface types")
		}
	}

	c.As = append(c.As, ifaces...)
	return nil
})

var Export = options.Option(func(c *Config, export bool) {
	c.Export = export
})

var Group = options.Option(func(c *Config, group string) {
	c.Group = strings.TrimSpace(group)
})

var LocationForPC = options.Option(func(c *Config, pc uintptr) {
	c.LocationPC = pc
})

var Name = options.Option(func(c *Config, name string) {
	c.Name = strings.TrimSpace(name)
})
