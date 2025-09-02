package decorator

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Decorator = (*decorator)(nil)

type Decorator interface {
	dependencies.Dependency
	Constructor() any
	Options() []Option
}

type decorator struct {
	constructor any
	options     []Option
}

func DecoratorFunc(constructor any, opts ...Option) Decorator {
	return &decorator{
		constructor,
		opts,
	}
}

func (d decorator) Constructor() any { return d.constructor }

func (d decorator) Options() []Option { return d.options }

func (d decorator) IsDependency() {}
