package dependencies

import (
	"frisboo-bank/pkg/options"
)

type Hooks interface {
	Dependency
	StartConstructor() any
	StopConstructor() any
	Options() *options.OptionBuilder[HooksOptions]
}

var _ Hooks = (*hooks)(nil)

type HooksOptions struct{}

func HookWithOptions() *options.OptionBuilder[HooksOptions] {
	return options.Apply(&HooksOptions{})
}

type hooks struct {
	startConstructor any
	stopConstructor  any
	options          *options.OptionBuilder[HooksOptions]
}

func Hook(
	startConstructor any,
	stopConstructor any,
	opts ...*options.OptionBuilder[HooksOptions],
) Hooks {
	var opt *options.OptionBuilder[HooksOptions]
	if len(opts) > 0 {
		opt = opts[0]
	}

	return &hooks{startConstructor, stopConstructor, opt}
}

func (h *hooks) StartConstructor() any { return h.startConstructor }

func (h *hooks) StopConstructor() any { return h.stopConstructor }

func (h *hooks) Options() *options.OptionBuilder[HooksOptions] { return h.options }

func (h *hooks) IsDependency() {}
