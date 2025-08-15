package dig

import (
	"context"
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/customerrors"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/utils"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

const (
	hookGroupStartPattern = "hook-start-%d"
	hookGroupStopPattern  = "hook-stop-%d"
)

var pError = customerrors.PrefixedError("dig container")

var _ contracts.ContainerAdapter = (*digAdapter)(nil)

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

func New(logger loggerContracts.Logger, waiter waiterContracts.Waiter) contracts.ContainerAdapter {
	utils.Assert(logger != nil, pError.New("logger can't be nil"))
	utils.Assert(waiter != nil, pError.New("waiter can't be nil"))

	return &digAdapter{
		dig:    dig.New(),
		logger: logger,
		waiter: waiter,
	}
}

func (d *digAdapter) RegisterDecorator(decorators ...decorator.Decorator) error {
	for _, decorator := range decorators {
		var opts []dig.DecorateOption
		if decorator.Options() != nil {
			opts = ToDigDecoratorOptions(*decorator.Options().Build())
		}

		if err := d.dig.Decorate(decorator.Constructor(), opts...); err != nil {
			return fmt.Errorf("failed to register decorator with error: %w", err)
		}
	}

	return nil
}

func (d *digAdapter) RegisterHook(hooks ...hook.Hooks) error {
	for _, hook := range hooks {
		groupID := len(d.hookGroups) + 1
		startGroup := fmt.Sprintf(hookGroupStartPattern, groupID)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, groupID)

		var opts []dig.ProvideOption
		if hook.Options() != nil {
			opts = ToDigHookOptions(*hook.Options().Build())
		}

		startOpts := append(opts, dig.Group(startGroup))
		if err := d.dig.Provide(hook.StartConstructor(), startOpts...); err != nil {
			return fmt.Errorf("failed to register hook start with error: %w", err)
		}

		stopOpts := append(opts, dig.Group(stopGroup))
		if err := d.dig.Provide(hook.StopConstructor(), stopOpts...); err != nil {
			return fmt.Errorf("failed to register hook stop with error: %w", err)
		}

		d.hookGroups = append(d.hookGroups, hookGroups{
			start: startGroup,
			stop:  stopGroup,
		})
	}

	return nil
}

func (d *digAdapter) RegisterInvoker(invokers ...invoker.Invoker) error {
	for _, invoker := range invokers {
		var opts []dig.InvokeOption
		if invoker.Options() != nil {
			opts = ToDigInvokerOptions(*invoker.Options().Build())
		}

		if err := d.dig.Invoke(invoker.Constructor(), opts...); err != nil {
			return fmt.Errorf("[dig] failed to register invoker with error: %w", err)
		}
	}

	return nil
}

func (d *digAdapter) RegisterProvider(providers ...provider.Provider) error {
	for _, provider := range providers {
		var opts []dig.ProvideOption
		if provider.Options() != nil {
			opts = ToDigProvideOptions(*provider.Options().Build())
		}

		if err := d.dig.Provide(provider.Constructor(), opts...); err != nil {
			return fmt.Errorf("failed to register provider with error: %w", err)
		}
	}

	return nil
}

func (d *digAdapter) Setup(cfg *config.Config) error {
	d.cfg = cfg

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
