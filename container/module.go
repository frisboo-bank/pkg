package container

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

func (m *module) apply(mp *module) {
	mp.modules = append(mp.modules, m)
}

func (m module) GetDecorators() []Decorator {
	return m.decorators
}

func (m module) GetHooks() []HookStarter {
	return m.hooks
}

func (m module) GetInvokes() []Invoker {
	return m.invokers
}

func (m module) GetModules() []Module {
	return m.modules
}

func (m module) GetName() string {
	return m.name
}

func (m module) GetProviders() []Provider {
	return m.providers
}
