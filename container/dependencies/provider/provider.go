package provider

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Provider = (*provider)(nil)

type Provider interface {
	dependencies.Dependency
	Fn() any
	Options() []Option
}

type provider struct {
	fn      any
	options []Option
}

func ProvideFunc(fn any, opts ...Option) Provider {
	return &provider{fn, opts}
}

func (i *provider) Fn() any { return i.fn }

func (i *provider) Options() []Option { return i.options }

func (i *provider) IsDependency() {}
