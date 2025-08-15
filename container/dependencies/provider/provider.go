package provider

import (
	"errors"
	"reflect"
	"strings"

	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/options"
)

type Provider interface {
	dependencies.Dependency
	Constructor() any
	Options() *options.OptionBuilder[ProviderOptions]
}

var _ Provider = (*provider)(nil)

type ProviderOptions struct {
	As         []any
	Export     bool
	Group      string
	LocationPC uintptr
	Name       string
}

func ProvideWithOptions() *options.OptionBuilder[ProviderOptions] {
	return options.Apply(&ProviderOptions{})
}

// As declares that the provided concrete type should also be exposed as the
// listed interface pointer types (pointers to interface types).
func As(ifaces ...any) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		for _, i := range ifaces {
			t := reflect.TypeOf(i)
			if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Interface {
				return errors.New("provider options: As expects pointers to interface types")
			}
		}

		opts.As = append(opts.As, ifaces...)
		return nil
	})
}

// Export marks this constructor as exported (adapter-dependent behavior).
func Export(export bool) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		opts.Export = export
		return nil
	})
}

func Group(group string) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		group = strings.TrimSpace(group)

		if group == "" {
			return errors.New("provider options: group can't be empty")
		}

		opts.Group = group
		return nil
	})
}

// LocationForPC sets an alternate PC for improved stack/introspection info.
func LocationForPC(pc uintptr) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		if pc == 0 {
			return errors.New("provider options: Location PC cannot be 0")
		}

		opts.LocationPC = pc
		return nil
	})
}

func Name(name string) options.Option[ProviderOptions] {
	return options.OptionFunc[ProviderOptions](func(opts *ProviderOptions) error {
		name = strings.TrimSpace(name)

		if name == "" {
			return errors.New("provider options: name can't be empty")
		}

		opts.Name = name
		return nil
	})
}

type provider struct {
	constructor any
	options     *options.OptionBuilder[ProviderOptions]
}

func Provide(constructor any, opts ...*options.OptionBuilder[ProviderOptions]) Provider {
	var opt *options.OptionBuilder[ProviderOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &provider{constructor, opt}
}

func (i *provider) Constructor() any { return i.constructor }

func (i *provider) Options() *options.OptionBuilder[ProviderOptions] { return i.options }

func (i *provider) IsDependency() {}
