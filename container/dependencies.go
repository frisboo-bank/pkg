package container

import "frisboo-bank/pkg/container/contracts"

var (
	_ contracts.Decorator   = (*decorator)(nil)
	_ contracts.HookStarter = (*hookStarter)(nil)
	_ contracts.Invoker     = (*invoker)(nil)
	_ contracts.Provider    = (*provider)(nil)
)

type decorator struct {
	fn      contracts.DecoratorFn
	options []contracts.DecoratorOption
}

type hookStarter struct {
	startFn contracts.HookStartFn
	stopFn  contracts.HookStopFn
	options []contracts.HookOption
}

type invoker struct {
	fn      contracts.InvokerFn
	options []contracts.InvokerOption
}

type provider struct {
	fn      contracts.ProviderFn
	options []contracts.ProviderOption
}

func Decorate(fn contracts.DecoratorFn, options ...contracts.DecoratorOption) contracts.Decorator {
	return &decorator{fn, options}
}
func (d *decorator) Fn() contracts.DecoratorFn            { return d.fn }
func (d *decorator) Options() []contracts.DecoratorOption { return d.options }
func (d *decorator) IsDependency()                        {}

func Hook(
	startFn contracts.HookStartFn,
	stopFn contracts.HookStopFn,
	options ...contracts.HookOption,
) contracts.HookStarter {
	return &hookStarter{startFn, stopFn, options}
}
func (h *hookStarter) Options() []contracts.HookOption { return h.options }
func (h *hookStarter) StartFn() contracts.HookStartFn  { return h.startFn }
func (h *hookStarter) StopFn() contracts.HookStopFn    { return h.stopFn }
func (h *hookStarter) IsDependency()                   {}

func Invoke(fn contracts.InvokerFn, options ...contracts.InvokerOption) contracts.Invoker {
	return &invoker{fn, options}
}
func (i *invoker) Fn() contracts.InvokerFn            { return i.fn }
func (i *invoker) Options() []contracts.InvokerOption { return i.options }
func (i *invoker) IsDependency()                      {}

func Provide(fn contracts.ProviderFn, options ...contracts.ProviderOption) contracts.Provider {
	return &provider{fn, options}
}
func (p *provider) Fn() contracts.ProviderFn            { return p.fn }
func (p *provider) Options() []contracts.ProviderOption { return p.options }
func (p *provider) IsDependency()                       {}
