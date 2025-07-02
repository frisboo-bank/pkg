package container

// Module represents a dependency-injectable module in the container system.
// Each module can have submodules, providers, decorators, invokes, and hooks.
type Module interface {
	Dependency
	GetName() string
	GetModules() []Module
	GetProviders() []Provider
	GetDecorators() []Decorator
	GetInvokes() []Invoker
	GetHooks() []HookStarter
}

type module struct {
	name       string
	modules    []Module
	providers  []Provider
	decorators []Decorator
	invokers   []Invoker
	hooks      []HookStarter
}

var _ Module = (*module)(nil)

func NewModule(name string, options ...Dependency) Module {
	m := &module{name: name}

	for _, option := range options {
		option.apply(m)
	}

	return m
}

// apply adds the module as a submodule to the target module pointer.
func (m *module) apply(targetModule *module) {
	for _, mod := range targetModule.modules {
		if mod == m {
			return
		}
	}

	targetModule.modules = append(targetModule.modules, m)
}

// returns the name of the module.
func (m *module) GetName() string {
	return m.name
}

// returns the submodules registered in the module.
func (m *module) GetModules() []Module {
	return m.modules
}

// returns the providers registered in the module.
func (m *module) GetProviders() []Provider {
	return m.providers
}

// returns the hooks registered in the module.
func (m *module) GetHooks() []HookStarter {
	return m.hooks
}

// returns the decorators registered in the module.
func (m *module) GetDecorators() []Decorator {
	return m.decorators
}

// returns the invokers registered in the module.
func (m *module) GetInvokes() []Invoker {
	return m.invokers
}
