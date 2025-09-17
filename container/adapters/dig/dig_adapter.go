package dig

import (
	"context"
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

const (
	hookGroupStartPattern = "hook-start-%d"
	hookGroupStopPattern  = "hook-stop-%d"
)

var _ contracts.ContainerAdapter = (*digAdapter)(nil)

type DigAdapterConfig struct{}

type hookGroups struct {
	start string
	stop  string
}

type digAdapter struct {
	cfg        *config.Config
	dig        *dig.Container
	logger     loggerContracts.Logger
	waiter     waiterContracts.Waiter
	hookGroups []hookGroups
}

func New(
	cfg *config.Config,
	waiter waiterContracts.Waiter,
	logger loggerContracts.Logger,
) contracts.ContainerAdapter {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("waiter", waiter)
	validation.AssertNotNil("logger", logger)

	return &digAdapter{
		cfg:    cfg,
		dig:    dig.New(),
		waiter: waiter,
		logger: logger,
	}
}

func (d *digAdapter) RegisterProvider(providers ...provider.Provider) error {
	for _, i := range providers {
		cfg := provider.Config{}
		if err := options.Apply(&cfg, i.Options()...); err != nil {
			return syserrors.Wrap(err, "failed to apply provider options")
		}
		opts := toDigProvideOptions(cfg)

		if err := d.dig.Provide(i.Constructor(), opts...); err != nil {
			return syserrors.Wrap(err, "failed to register provider")
		}
	}
	return nil
}

func (d *digAdapter) RegisterDecorator(decorators ...decorator.Decorator) error {
	for _, i := range decorators {
		cfg := decorator.Config{}
		if err := options.Apply(&cfg, i.Options()...); err != nil {
			return syserrors.Wrap(err, "failed to apply decorator options")
		}
		opts := toDigDecoratorOptions(cfg)

		if err := d.dig.Decorate(i.Constructor(), opts...); err != nil {
			return syserrors.Wrap(err, "failed to register decorator")
		}
	}
	return nil
}

func (d *digAdapter) RegisterHook(hooks ...hook.Hooks) error {
	for _, i := range hooks {
		groupID := len(d.hookGroups) + 1
		startGroup := fmt.Sprintf(hookGroupStartPattern, groupID)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, groupID)

		cfg := hook.Config{}
		if err := options.Apply(&cfg, i.Options()...); err != nil {
			return syserrors.Wrap(err, "failed to apply hook options")
		}
		opts := toDigHookOptions(cfg)

		startOpts := append(opts, dig.Group(startGroup))
		if err := d.dig.Provide(i.StartConstructor(), startOpts...); err != nil {
			return syserrors.Wrap(err, "failed to register hook start")
		}

		stopOpts := append(opts, dig.Group(stopGroup))
		if err := d.dig.Provide(i.StopConstructor(), stopOpts...); err != nil {
			return syserrors.Wrap(err, "failed to register hook stop")
		}

		d.hookGroups = append(d.hookGroups, hookGroups{
			start: startGroup,
			stop:  stopGroup,
		})
	}
	return nil
}

func (d *digAdapter) RegisterInvoker(invokers ...invoker.Invoker) error {
	for _, i := range invokers {
		cfg := invoker.Config{}
		if err := options.Apply(&cfg, i.Options()...); err != nil {
			return syserrors.Wrap(err, "failed to apply invoker options")
		}
		opts := toDigInvokerOptions(cfg)

		if err := d.dig.Invoke(i.Constructor(), opts...); err != nil {
			return syserrors.Wrap(err, "failed to register invoker")
		}
	}
	return nil
}

func (d *digAdapter) Start(_ context.Context) error {
	hooks, err := d.resolveHooks()
	if err != nil {
		return syserrors.Wrap(err, "failed to resolve hook")
	}

	d.waiter.Add(hooks...)

	return d.waiter.Wait()
}

func (d *digAdapter) Stop(_ context.Context) error {
	d.waiter.Cancel()
	return nil
}

func (d *digAdapter) Type() containertype.ContainerType {
	return containertype.ContainerTypes.DIG
}

func (d *digAdapter) resolveHooks() ([]waiterContracts.WaiterHook, error) {
	hooks := make([]waiterContracts.WaiterHook, len(d.hookGroups))

	for i, hookGroup := range d.hookGroups {
		hook := waiterContracts.WaiterHook{}

		startFns, err := resolveDynamicGroup[[]waiterContracts.WaitFunc](
			d.dig,
			hookGroup.start,
			reflect.TypeOf([]waiterContracts.WaitFunc{}),
		)
		if err != nil {
			return nil, err
		}

		if len(startFns) > 0 {
			hook.Wait = startFns[0]
		}

		stopFns, err := resolveDynamicGroup[[]waiterContracts.CleanupFunc](
			d.dig,
			hookGroup.stop,
			reflect.TypeOf([]waiterContracts.CleanupFunc{}),
		)
		if err != nil {
			return nil, err
		}

		if len(stopFns) > 0 {
			hook.Cleanup = stopFns[0]
		}

		hooks[i] = hook
	}

	return hooks, nil
}
