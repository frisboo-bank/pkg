package invoker

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Invoker = (*invoker)(nil)

type Invoker interface {
	dependencies.Dependency
	Constructor() any
	Options() []Option
}

type invoker struct {
	constructor any
	options     []Option
}

func InvokerFunc(constructor any, opts ...Option) Invoker {
	return &invoker{
		constructor,
		opts,
	}
}

func (i *invoker) Constructor() any { return i.constructor }

func (i *invoker) Options() []Option { return i.options }

func (i *invoker) IsDependency() {}
