package container

import "frisboo-bank/pkg/container/contracts"

var _ contracts.Module = (*module)(nil)

type module struct {
	decorators []contracts.Decorator
	hooks      []contracts.HookStarter
	invokers   []contracts.Invoker
	modules    []contracts.Module
	name       string
	providers  []contracts.Provider
}

func NewModule(name string, deps ...contracts.Dependency) contracts.Module {
	m := &module{name: name}

	for _, dep := range deps {
		switch d := dep.(type) {
		case contracts.Decorator:
			m.decorators = append(m.decorators, d)
		case contracts.HookStarter:
			m.hooks = append(m.hooks, d)
		case contracts.Invoker:
			m.invokers = append(m.invokers, d)
		case contracts.Module:
			m.modules = append(m.modules, d)
		case contracts.Provider:
			m.providers = append(m.providers, d)
		}
	}

	return m
}

func (m *module) Decorators() []contracts.Decorator { return m.decorators }
func (m *module) Hooks() []contracts.HookStarter    { return m.hooks }
func (m *module) Invokers() []contracts.Invoker     { return m.invokers }
func (m *module) Modules() []contracts.Module       { return m.modules }
func (m *module) Name() string                      { return m.name }
func (m *module) Providers() []contracts.Provider   { return m.providers }
func (m *module) IsDependency()                     {}
