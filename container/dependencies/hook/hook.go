package hook

import (
	"frisboo-bank/pkg/container/dependencies"
)

var _ Hooks = (*hooks)(nil)

type Hooks interface {
	dependencies.Dependency

	Name() string
	StartFn() any
	StopFn() any
	Options() []Option
}

type hooks struct {
	name    string
	startFn any
	stopFn  any
	options []Option
}

func HooksFunc(name string, startFn any, stopFn any, opts ...Option) Hooks {
	return &hooks{
		name,
		startFn,
		stopFn,
		opts,
	}
}

func (h *hooks) Name() string      { return h.name }
func (h *hooks) StartFn() any      { return h.startFn }
func (h *hooks) StopFn() any       { return h.stopFn }
func (h *hooks) Options() []Option { return h.options }
func (h *hooks) IsDependency()     {}
