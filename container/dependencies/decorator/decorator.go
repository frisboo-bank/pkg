package decorator

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Decorator = (*decorator)(nil)

type Decorator interface {
	dependencies.Dependency
	Fn() any
	Options() []Option
}

type decorator struct {
	fn      any
	options []Option
}

func DecoratorFunc(fn any, opts ...Option) Decorator {
	return &decorator{fn, opts}
}

func (d *decorator) Fn() any           { return d.fn }
func (d *decorator) Options() []Option { return d.options }
func (d *decorator) IsDependency()     {}
