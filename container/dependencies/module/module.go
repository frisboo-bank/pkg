package module

import (
	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
)

type Module interface {
	dependencies.Dependency
	Decorators() []decorator.Decorator
	Hooks() []hook.Hooks
	Invokers() []invoker.Invoker
	Modules() []Module
	Name() string
	Providers() []provider.Provider
}

var _ Module = (*module)(nil)

type module struct {
	decorators []decorator.Decorator
	hooks      []hook.Hooks
	invokers   []invoker.Invoker
	modules    []Module
	name       string
	providers  []provider.Provider
}

func NewModule(name string, deps ...dependencies.Dependency) Module {
	m := &module{name: name}

	for _, dep := range deps {
		switch d := dep.(type) {
		case decorator.Decorator:
			m.decorators = append(m.decorators, d)
		case hook.Hooks:
			m.hooks = append(m.hooks, d)
		case invoker.Invoker:
			m.invokers = append(m.invokers, d)
		case Module:
			m.modules = append(m.modules, d)
		case provider.Provider:
			m.providers = append(m.providers, d)
		}
	}

	return m
}

func (m *module) Decorators() []decorator.Decorator { return m.decorators }

func (m *module) Hooks() []hook.Hooks { return m.hooks }

func (m *module) Invokers() []invoker.Invoker { return m.invokers }

func (m *module) IsDependency() {}

func (m *module) Modules() []Module { return m.modules }

func (m *module) Name() string { return m.name }

func (m *module) Providers() []provider.Provider { return m.providers }
