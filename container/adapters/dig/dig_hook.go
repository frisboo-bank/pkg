package dig

import (
	"reflect"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

func (a *digAdapter) RegisterHooks(hooks ...hook.Hooks) error {
	for _, h := range hooks {
		if err := a.RegisterHook(h); err != nil {
			return err
		}
	}
	return nil
}

func (a *digAdapter) RegisterHook(h hook.Hooks) error {
	name := h.Name()

	cfg := hook.Config{}
	if err := options.Apply(&cfg, h.Options()...); err != nil {
		return syserrors.Wrapf(err, "failed to apply hook %s options", name)
	}

	startFunc := h.StartFn()
	stopFunc := h.StopFn()

	if startFunc == nil && stopFunc == nil {
		return nil
	}

	var err error
	if startFunc != nil {
		startFunc, err = wrapFuncWithDigIn(startFunc, cfg.NamedDeps, "hook "+name+" start")
		if err != nil {
			return err
		}
	}
	if stopFunc != nil {
		stopFunc, err = wrapFuncWithDigIn(stopFunc, cfg.NamedDeps, "hook "+name+" stop")
		if err != nil {
			return err
		}
	}

	opts := toDigHookOptions(cfg)

	startGroup := name + "-start"
	if startFunc != nil {
		if err := a.dig.Provide(startFunc, append(opts, dig.Group(startGroup))...); err != nil {
			return syserrors.Wrapf(err, "failed to register hook start %s", startGroup)
		}
	}
	stopGroup := name + "-stop"
	if stopFunc != nil {
		if err := a.dig.Provide(stopFunc, append(opts, dig.Group(stopGroup))...); err != nil {
			return syserrors.Wrapf(err, "failed to register hook stop %s", stopGroup)
		}
	}

	a.hooks = append(a.hooks, name)

	return nil
}

func (a *digAdapter) resolveHooks() ([]waiterContracts.WaiterHook, error) {
	hooks := make([]waiterContracts.WaiterHook, len(a.hooks))
	for i, h := range a.hooks {
		h, err := a.resolveHook(h)
		if err != nil {
			return nil, err
		}
		hooks[i] = h
	}
	return hooks, nil
}

func (a *digAdapter) resolveHook(hook string) (waiterContracts.WaiterHook, error) {
	wh := waiterContracts.WaiterHook{Name: hook}

	startGroup := hook + "-start"
	startFns, err := resolveDynamicGroup[[]waiterContracts.WaitFunc](
		a.dig,
		startGroup,
		reflect.TypeOf([]waiterContracts.WaitFunc{}),
	)
	if err != nil {
		return wh, syserrors.Wrapf(err, "failed to resolve start hook %s", startGroup)
	}
	if len(startFns) > 0 {
		wh.Wait = startFns[0]
	}

	stopGroup := hook + "-stop"
	stopFns, err := resolveDynamicGroup[[]waiterContracts.CleanupFunc](
		a.dig,
		stopGroup,
		reflect.TypeOf([]waiterContracts.CleanupFunc{}),
	)
	if err != nil {
		return wh, syserrors.Wrapf(err, "failed to resolve stop hook %s", stopGroup)
	}
	if len(stopFns) > 0 {
		wh.Cleanup = stopFns[0]
	}

	return wh, nil
}

func toDigHookOptions(cfg hook.Config) []dig.ProvideOption {
	var result []dig.ProvideOption

	if len(cfg.As) > 0 {
		result = append(result, dig.As(cfg.As...))
	}

	if cfg.Export {
		result = append(result, dig.Export(true))
	}

	if cfg.Group != "" {
		result = append(result, dig.Group(cfg.Group))
	}

	if cfg.LocationPC != 0 {
		result = append(result, dig.LocationForPC(cfg.LocationPC))
	}

	if cfg.Name != "" {
		result = append(result, dig.Name(cfg.Name))
	}

	return result
}
