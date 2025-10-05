package invoker

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Invoker = (*invoker)(nil)

type Invoker interface {
	dependencies.Dependency
	Fn() any
	Options() []Option
}

type invoker struct {
	fn      any
	options []Option
}

func InvokerFunc(fn any, opts ...Option) Invoker {
	return &invoker{fn, opts}
}

func (i *invoker) Fn() any           { return i.fn }
func (i *invoker) Options() []Option { return i.options }
func (i *invoker) IsDependency()     {}
