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

type DecoratorOptions struct{}

func DecorateWithOptions() *options.OptionBuilder[DecoratorOptions] {
	return options.Apply(&DecoratorOptions{})
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
