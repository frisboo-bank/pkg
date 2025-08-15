package hook

import (
	"errors"
	"reflect"
	"strings"

	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/options"
)

type Hooks interface {
	dependencies.Dependency
	StartConstructor() any
	StopConstructor() any
	Options() *options.OptionBuilder[HooksOptions]
}

var _ Hooks = (*hooks)(nil)

type HooksOptions struct {
	As         []any
	Export     bool
	Group      string
	LocationPC uintptr
	Name       string
}

func HookWithOptions() *options.OptionBuilder[HooksOptions] {
	return options.Apply(&HooksOptions{})
}

// As declares that the provided concrete type should also be exposed as the
// listed interface pointer types (pointers to interface types).
func As(ifaces ...any) options.Option[HooksOptions] {
	return options.OptionFunc[HooksOptions](func(opts *HooksOptions) error {
		for _, i := range ifaces {
			t := reflect.TypeOf(i)
			if t == nil || t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Interface {
				return errors.New("hook options: As expects pointers to interface types")
			}
		}

		opts.As = append(opts.As, ifaces...)
		return nil
	})
}

// Export marks this constructor as exported (adapter-dependent behavior).
func Export(export bool) options.Option[HooksOptions] {
	return options.OptionFunc[HooksOptions](func(opts *HooksOptions) error {
		opts.Export = export
		return nil
	})
}

func Group(group string) options.Option[HooksOptions] {
	return options.OptionFunc[HooksOptions](func(opts *HooksOptions) error {
		group = strings.TrimSpace(group)

		if group == "" {
			return errors.New("hook options: group can't be empty")
		}

		opts.Group = group
		return nil
	})
}

// LocationForPC sets an alternate PC for improved stack/introspection info.
func LocationForPC(pc uintptr) options.Option[HooksOptions] {
	return options.OptionFunc[HooksOptions](func(opts *HooksOptions) error {
		if pc == 0 {
			return errors.New("hook options: Location PC cannot be 0")
		}

		opts.LocationPC = pc
		return nil
	})
}

func Name(name string) options.Option[HooksOptions] {
	return options.OptionFunc[HooksOptions](func(opts *HooksOptions) error {
		name = strings.TrimSpace(name)

		if name == "" {
			return errors.New("provider options: name can't be empty")
		}

		opts.Name = name
		return nil
	})
}

type hooks struct {
	startConstructor any
	stopConstructor  any
	options          *options.OptionBuilder[HooksOptions]
}

func Hook(
	startConstructor any,
	stopConstructor any,
	opts ...*options.OptionBuilder[HooksOptions],
) Hooks {
	var opt *options.OptionBuilder[HooksOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &hooks{startConstructor, stopConstructor, opt}
}

func (h *hooks) StartConstructor() any { return h.startConstructor }

func (h *hooks) StopConstructor() any { return h.stopConstructor }

func (h *hooks) Options() *options.OptionBuilder[HooksOptions] { return h.options }

func (h *hooks) IsDependency() {}
