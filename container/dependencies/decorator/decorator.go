package decorator

import (
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/options"
)

type Decorator interface {
	dependencies.Dependency
	Constructor() any
	Options() *options.OptionBuilder[DecoratorOptions]
}

var _ Decorator = (*decorator)(nil)

type DecoratorOptions struct {
	BeforeCallback any
	Callback       any
	Info           any
}

func DecorateWithOptions() *options.OptionBuilder[DecoratorOptions] {
	return options.Apply(&DecoratorOptions{})
}

// WithBeforeCallback registers a pre-execution callback (e.g. dig.BeforeCallback).
func WithBeforeCallback(cb any) options.Option[DecoratorOptions] {
	return options.OptionFunc[DecoratorOptions](func(o *DecoratorOptions) error {
		o.BeforeCallback = cb
		return nil
	})
}

// WithCallback registers a post-execution callback (e.g. dig.Callback).
func WithCallback(cb any) options.Option[DecoratorOptions] {
	return options.OptionFunc[DecoratorOptions](func(o *DecoratorOptions) error {
		o.Callback = cb
		return nil
	})
}

// WithInfo sets a pointer to a structure that will be filled with
// framework-specific decorate info (e.g. *dig.DecorateInfo).
func WithInfo(info any) options.Option[DecoratorOptions] {
	return options.OptionFunc[DecoratorOptions](func(o *DecoratorOptions) error {
		o.Info = info
		return nil
	})
}

type decorator struct {
	constructor any
	options     *options.OptionBuilder[DecoratorOptions]
}

func Decorate(constructor any, opts ...*options.OptionBuilder[DecoratorOptions]) Decorator {
	var opt *options.OptionBuilder[DecoratorOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &decorator{constructor, opt}
}

func (d *decorator) Constructor() any { return d.constructor }

func (d *decorator) Options() *options.OptionBuilder[DecoratorOptions] { return d.options }

func (d *decorator) IsDependency() {}
