package container

// Dependency is an interface for anything that can be applied to a module.
type Dependency interface {
	apply(*module)
}

type (
	ProviderFn     any
	ProviderOption any
)

type Provider interface {
	Dependency
	Fn() ProviderFn
	Options() []ProviderOption
}

type provider struct {
	fn      ProviderFn
	options []ProviderOption
}

func (d provider) Fn() ProviderFn {
	return d.fn
}

func (d provider) Options() []ProviderOption {
	return d.options
}

func (d provider) apply(module *module) {
	module.providers = append(module.providers, d)
}

func Provide(fn ProviderFn, options ...ProviderOption) Provider {
	return provider{fn, options}
}

type (
	DecoratorFn     any
	DecoratorOption any
)

type Decorator interface {
	Dependency
	Fn() DecoratorFn
	Options() []DecoratorOption
}

type decorator struct {
	fn      DecoratorFn
	options []DecoratorOption
}

func (d decorator) Fn() DecoratorFn {
	return d.fn
}

func (d decorator) Options() []DecoratorOption {
	return d.options
}

func (d decorator) apply(module *module) {
	module.decorators = append(module.decorators, d)
}

func Decorate(fn DecoratorFn, options ...DecoratorOption) Decorator {
	return decorator{fn, options}
}

type (
	InvokerFn     any
	InvokerOption any
)

type Invoker interface {
	Dependency
	Fn() InvokerFn
	Options() []InvokerOption
}

type invoker struct {
	fn      InvokerFn
	options []InvokerOption
}

func (d invoker) Fn() InvokerFn {
	return d.fn
}

// Options implements Invoker.
func (d invoker) Options() []InvokerOption {
	return d.options
}

func (o invoker) apply(module *module) {
	module.invokers = append(module.invokers, o)
}

func Invoke(fn InvokerFn, options ...InvokerOption) Invoker {
	return invoker{fn, options}
}

type (
	HookStartFn any
	HookStopFn  any
	HookOption  any
)

type HookStarter interface {
	Dependency
	StartFn() HookStartFn
	StopFn() HookStopFn
	Options() []HookOption
}

type hookStarter struct {
	startFn HookStartFn
	stopFn  HookStopFn
	options []HookOption
}

func (d hookStarter) StartFn() HookStartFn {
	return d.startFn
}

func (d hookStarter) StopFn() HookStopFn {
	return d.stopFn
}

func (d hookStarter) Options() []HookOption {
	return d.options
}

func (d hookStarter) apply(module *module) {
	module.hooks = append(module.hooks, d)
}

func Hook(startFn HookStartFn, stopFn HookStopFn, options ...HookOption) HookStarter {
	return hookStarter{startFn, stopFn, options}
}
