package dig

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

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
	syserrors.AssertNotNil("cfg", cfg)
	syserrors.AssertNotNil("waiter", waiter)
	syserrors.AssertNotNil("logger", logger)

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
			return err
		}
		opts := toDigProvideOptions(cfg)

		if err := d.dig.Provide(i.Constructor(), opts...); err != nil {
			if verr := dig.Visualize(d.dig, os.Stdout, dig.VisualizeError(err)); verr != nil {
				return err
			}
			return syserrors.Newf("failed to register provider with error: %w", err)
		}
	}
	return nil
}

func (d *digAdapter) RegisterDecorator(decorators ...decorator.Decorator) error {
	for _, i := range decorators {
		cfg := decorator.Config{}
		if err := options.Apply(&cfg, i.Options()...); err != nil {
			return err
		}
		opts := toDigDecoratorOptions(cfg)

		if err := d.dig.Decorate(i.Constructor(), opts...); err != nil {
			return syserrors.Newf("failed to register decorator with error: %w", err)
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
			return err
		}
		opts := toDigHookOptions(cfg)

		startOpts := append(opts, dig.Group(startGroup))
		if err := d.dig.Provide(i.StartConstructor(), startOpts...); err != nil {
			return syserrors.Newf("failed to register hook start with error: %w", err)
		}

		stopOpts := append(opts, dig.Group(stopGroup))
		if err := d.dig.Provide(i.StopConstructor(), stopOpts...); err != nil {
			return syserrors.Newf("failed to register hook stop with error: %w", err)
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
			return err
		}
		opts := toDigInvokerOptions(cfg)

		if err := d.dig.Invoke(i.Constructor(), opts...); err != nil {
			return syserrors.Newf("[dig] failed to register invoker with error: %w", err)
		}
	}
	return nil
}

func (d *digAdapter) Start(_ context.Context) error {
	hooks, err := d.resolveHooks()
	if err != nil {
		return err
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
