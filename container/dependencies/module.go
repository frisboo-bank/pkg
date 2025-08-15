package dependencies

type Dependency interface {
	IsDependency()
}

type Module interface {
	Dependency
	Decorators() []Decorator
	Hooks() []Hooks
	Invokers() []Invoker
	Modules() []Module
	Name() string
	Providers() []Provider
}

var _ Module = (*module)(nil)

type module struct {
	decorators []Decorator
	hooks      []Hooks
	invokers   []Invoker
	modules    []Module
	name       string
	providers  []Provider
}

func NewModule(name string, deps ...Dependency) Module {
	m := &module{name: name}

	for _, dep := range deps {
		switch d := dep.(type) {
		case Decorator:
			m.decorators = append(m.decorators, d)
		case Hooks:
			m.hooks = append(m.hooks, d)
		case Invoker:
			m.invokers = append(m.invokers, d)
		case Module:
			m.modules = append(m.modules, d)
		case Provider:
			m.providers = append(m.providers, d)
		}
	}

	return m
}

func (m *module) Decorators() []Decorator { return m.decorators }

func (m *module) Hooks() []Hooks { return m.hooks }

func (m *module) Invokers() []Invoker { return m.invokers }

func (m *module) Modules() []Module { return m.modules }

func (m *module) Name() string { return m.name }

func (m *module) Providers() []Provider { return m.providers }

func (m *module) IsDependency() {}
