package dependencies

import (
	"frisboo-bank/pkg/options"
)

type Invoker interface {
	Dependency
	Constructor() any
	Options() *options.OptionBuilder[InvokerOptions]
}

var _ Invoker = (*invoker)(nil)

type InvokerOptions struct{}

func InvokeWithOptions() *options.OptionBuilder[InvokerOptions] {
	return options.Apply(&InvokerOptions{})
}

type invoker struct {
	constructor any
	options     *options.OptionBuilder[InvokerOptions]
}

func Invoke(constructor any, opts ...*options.OptionBuilder[InvokerOptions]) Invoker {
	var opt *options.OptionBuilder[InvokerOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &invoker{constructor, opt}
}

func (i *invoker) Constructor() any { return i.constructor }

func (i *invoker) Options() *options.OptionBuilder[InvokerOptions] { return i.options }

func (i *invoker) IsDependency() {}
