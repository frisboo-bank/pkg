package dig

import (
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"go.uber.org/dig"
)

func (a *digAdapter) RegisterProviders(providers ...provider.Provider) error {
	for id, p := range providers {
		name := fmt.Sprintf("provider-%d", id)
		if err := a.RegisterProvider(name, p); err != nil {
			return err
		}
	}
	return nil
}

func (a *digAdapter) RegisterProvider(name string, p provider.Provider) error {
	cfg := provider.Config{}
	if err := options.Apply(&cfg, p.Options()...); err != nil {
		return syserrors.Wrap(err, "failed to apply provider options")
	}

	fn := p.Fn()
	if len(cfg.NamedDeps) > 0 {
		var err error
		fn, err = wrapFuncWithNamedInputs(fn, cfg.NamedDeps, "provider "+name)
		if err != nil {
			return syserrors.Wrapf(err, "failed to adapt provider %s named deps", name)
		}
	}

	if cfg.Name == "" || cfg.Group == "" {
		if err := a.dig.Provide(fn, toDigProvideOptions(cfg)...); err != nil {
			return syserrors.Wrap(err, "failed to register provider")
		}
		return nil
	}

	// case we want to register group and name at the same time
	group := cfg.Group
	cfg.Group = ""

	if err := a.dig.Provide(fn, toDigProvideOptions(cfg)...); err != nil {
		return syserrors.Wrap(err, "failed to register provider")
	}

	// Reflect the function to find the primary return type.
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return syserrors.Newf("fn is not a function (kind=%s)", t.Kind())
	}
	numOut := t.NumOut()
	if numOut == 0 || numOut > 2 {
		return syserrors.Newf("fn must return (T) or (T, error); got %d outputs", numOut)
	}
	outType := t.Out(0)
	if numOut == 2 {
		second := t.Out(1)
		if !second.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			return syserrors.Newf("second return value is not error but %s", second.String())
		}
	}

	// Build the dig.In struct type that extracts the named instance.
	inParams := reflect.StructOf([]reflect.StructField{
		{
			Name:      "In",
			Type:      reflect.TypeOf(dig.In{}),
			Anonymous: true,
		},
		{
			Name: "Value",
			Type: outType,
			Tag:  reflect.StructTag(`name:"` + cfg.Name + `"`),
		},
	})

	fnType := reflect.FuncOf([]reflect.Type{inParams}, []reflect.Type{outType}, false)
	fnVal := reflect.MakeFunc(fnType, func(args []reflect.Value) []reflect.Value {
		in := args[0]
		return []reflect.Value{in.FieldByName("Value")}
	})

	// Re-export into group.
	cfg.Name = ""
	cfg.Group = group
	if err := a.dig.Provide(fnVal.Interface(), toDigProvideOptions(cfg)...); err != nil {
		return syserrors.Wrap(err, "failed to register provider")
	}

	return nil
}

func toDigProvideOptions(cfg provider.Config) []dig.ProvideOption {
	var result []dig.ProvideOption

	if len(cfg.As) > 0 {
		result = append(result, dig.As(cfg.As...))
	}

	if cfg.Export {
		result = append(result, dig.Export(true))
	}

	if cfg.Group != "" {
		result = append(result, dig.Group(cfg.Group))
	}

	if cfg.LocationPC != 0 {
		result = append(result, dig.LocationForPC(cfg.LocationPC))
	}

	if cfg.Name != "" {
		result = append(result, dig.Name(cfg.Name))
	}

	return result
}
