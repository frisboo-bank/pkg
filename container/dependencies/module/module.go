package module

import (
	"sync"

	"frisboo-bank/pkg/container/dependencies"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
)

type Module interface {
	dependencies.Dependency
	AddModules(modules ...Module) Module
	Modules() []Module
	AddProviders(providers ...provider.Provider) Module
	Providers() []provider.Provider
	AddDecorators(decorators ...decorator.Decorator) Module
	Decorators() []decorator.Decorator
	AddHooks(hooks ...hook.Hooks) Module
	Hooks() []hook.Hooks
	AddInvokers(invokers ...invoker.Invoker) Module
	Invokers() []invoker.Invoker
	Name() string
}

var _ Module = (*module)(nil)

type module struct {
	modules    []Module
	providers  []provider.Provider
	decorators []decorator.Decorator
	hooks      []hook.Hooks
	invokers   []invoker.Invoker
	name       string
	mu         sync.RWMutex
}

func ModuleFunc(name string, deps ...dependencies.Dependency) Module {
	m := &module{name: name}

	for _, dep := range deps {
		switch d := dep.(type) {
		case Module:
			m.AddModules(d)
		case provider.Provider:
			m.AddProviders(d)
		case decorator.Decorator:
			m.AddDecorators(d)
		case hook.Hooks:
			m.AddHooks(d)
		case invoker.Invoker:
			m.AddInvokers(d)
		}
	}

	return m
}

func (m *module) AddModules(modules ...Module) Module {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.modules = append(m.modules, modules...)
	return m
}

func (m *module) Modules() []Module {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.modules
}

func (m *module) AddProviders(providers ...provider.Provider) Module {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.providers = append(m.providers, providers...)
	return m
}

func (m *module) Providers() []provider.Provider {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.providers
}

func (m *module) AddDecorators(decorators ...decorator.Decorator) Module {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.decorators = append(m.decorators, decorators...)
	return m
}

func (m *module) Decorators() []decorator.Decorator {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.decorators
}

func (m *module) AddHooks(hooks ...hook.Hooks) Module {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hooks = append(m.hooks, hooks...)
	return m
}

func (m *module) Hooks() []hook.Hooks {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.hooks
}

func (m *module) AddInvokers(invokers ...invoker.Invoker) Module {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.invokers = append(m.invokers, invokers...)
	return m
}

func (m *module) Invokers() []invoker.Invoker {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.invokers
}

func (m *module) IsDependency() {}

func (m *module) Name() string { return m.name }
