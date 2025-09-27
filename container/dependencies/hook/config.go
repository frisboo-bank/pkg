package hook

import (
	"reflect"
	"strings"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	As         []any
	Export     bool
	Group      string
	LocationPC uintptr
	Name       string
	NamedDeps  map[string]string
}

type Option = options.OptionFn[Config]

func (c *Config) Validate() error {
	return validation.ValidateStruct(c)
	// validation.Field(&c.NamedDeps, validation.Map(
	// 	validation.Key("Iface", validation.Required),
	// 	validation.Key("Name", validation.Required),
	// )))
}

var As = options.VarOptionErr(func(c *Config, ifaces ...any) error {
	if len(ifaces) == 0 {
		return syserrors.CantBeEmptyError("As")
	}
	for _, i := range ifaces {
		t := reflect.TypeOf(i)
		if t == nil || t.Kind() != reflect.Pointer || t.Elem().Kind() != reflect.Interface {
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

func NamedDep(ref string, name string) Option {
	return func(c *Config) error {
		if c.NamedDeps == nil {
			c.NamedDeps = make(map[string]string, 1)
		}
		c.NamedDeps[ref] = name
		return nil
	}
}
