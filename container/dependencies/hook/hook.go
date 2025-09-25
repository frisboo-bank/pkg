package hook

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Hooks = (*hooks)(nil)

type Hooks interface {
	dependencies.Dependency
	StartConstructor() any
	StopConstructor() any
	Options() []Option
}

type hooks struct {
	startConstructor any
	stopConstructor  any
	options          []Option
}

func HooksFunc(startConstructor any, stopConstructor any, opts ...Option) Hooks {
	return &hooks{
		startConstructor,
		stopConstructor,
		opts,
	}
}

func (h *hooks) StartConstructor() any { return h.startConstructor }
func (h *hooks) StopConstructor() any  { return h.stopConstructor }
func (h *hooks) Options() []Option     { return h.options }
func (h *hooks) IsDependency()         {}
