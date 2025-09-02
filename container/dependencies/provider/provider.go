package provider

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Provider = (*provider)(nil)

type Provider interface {
	dependencies.Dependency
	Constructor() any
	Options() []Option
}

type provider struct {
	constructor any
	options     []Option
}

func ProvideFunc(constructor any, opts ...Option) Provider {
	return &provider{
		constructor,
		opts,
	}
}

func (i provider) Constructor() any { return i.constructor }

func (i provider) Options() []Option { return i.options }

func (i provider) IsDependency() {}
