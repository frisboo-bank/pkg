package invoker

import (
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/options"
)

type Invoker interface {
	dependencies.Dependency
	Constructor() any
	Options() *options.OptionBuilder[InvokerOptions]
}

var _ Invoker = (*invoker)(nil)

type InvokerOptions struct {
	Info any
}

func InvokeWithOptions() *options.OptionBuilder[InvokerOptions] {
	return options.Apply(&InvokerOptions{})
}

// WithInfo sets a pointer that will receive invocation info (e.g. *dig.InvokeInfo).
func WithInfo(info any) options.Option[InvokerOptions] {
	return options.OptionFunc[InvokerOptions](func(o *InvokerOptions) error {
		o.Info = info
		return nil
	})
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
