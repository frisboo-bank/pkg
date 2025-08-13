package contracts

type (
	ProviderFn     any
	ProviderOption any

	DecoratorFn     any
	DecoratorOption any

	InvokerFn     any
	InvokerOption any

	HookStartFn any
	HookStopFn  any
	HookOption  any
)

type Dependency interface {
	IsDependency()
}

type Module interface {
	Dependency
	Name() string
	Modules() []Module
	Providers() []Provider
	Decorators() []Decorator
	Invokers() []Invoker
	Hooks() []HookStarter
}

type Provider interface {
	Dependency
	Fn() ProviderFn
	Options() []ProviderOption
}

type Decorator interface {
	Dependency
	Fn() DecoratorFn
	Options() []DecoratorOption
}

type Invoker interface {
	Dependency
	Fn() InvokerFn
	Options() []InvokerOption
}

type HookStarter interface {
	Dependency
	StartFn() HookStartFn
	StopFn() HookStopFn
	Options() []HookOption
}
