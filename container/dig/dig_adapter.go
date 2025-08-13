package dig

import (
	"context"
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/utils"

	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

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

func (d *digAdapter) RegisterDecorator(decorators ...contracts.Decorator) error {
	for _, decorator := range decorators {
		opts, err := filterOptions[dig.DecorateOption](decorator.Options())
		if err != nil {
			return fmt.Errorf("failed to convert decorator options with error: %w", err)
		}

		if err := d.dig.Decorate(decorator.Fn(), opts...); err != nil {
			return fmt.Errorf("failed to register decorator with error: %w", err)
		}
	}

	return nil
}

func (d *digAdapter) RegisterHook(hooks ...contracts.HookStarter) error {
	for _, hook := range hooks {
		groudID := len(d.hookGroups) + 1
		startGroup := fmt.Sprintf(hookGroupStartPattern, groudID)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, groudID)

		opts, err := filterOptions[dig.ProvideOption](hook.Options())
		if err != nil {
			return fmt.Errorf("failed to convert hook options with error: %w", err)
		}

		startOpts := append(opts, dig.Group(startGroup))
		if err := d.dig.Provide(hook.StartFn(), startOpts...); err != nil {
			return fmt.Errorf("failed to register hook start with error: %w", err)
		}

		stopOpts := append(opts, dig.Group(stopGroup))
		if err := d.dig.Provide(hook.StopFn(), stopOpts...); err != nil {
			return fmt.Errorf("failed to register hook stop with error: %w", err)
		}

		d.hookGroups = append(d.hookGroups, hookGroups{
			start: startGroup,
			stop:  stopGroup,
		})
	}

	return nil
}

func (d *digAdapter) RegisterInvoker(invokers ...contracts.Invoker) error {
	for _, invoker := range invokers {
		opts, err := filterOptions[dig.InvokeOption](invoker.Options())
		if err != nil {
			return fmt.Errorf("[dig] failed to convert invoker options with error: %w", err)
		}

		if err := d.dig.Invoke(invoker.Fn(), opts...); err != nil {
			return fmt.Errorf("[dig] failed to register invoker with error: %w", err)
		}
	}

	return nil
}

func (d *digAdapter) RegisterProvider(providers ...contracts.Provider) error {
	for _, provider := range providers {
		opts, err := filterOptions[dig.ProvideOption](provider.Options())
		if err != nil {
			return fmt.Errorf("failed to convert provider options with error: %w", err)
		}

		if err := d.dig.Provide(provider.Fn(), opts...); err != nil {
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

	for _, hookGroup := range d.hookGroups {
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

		hooks = append(hooks, hook)
	}

	return hooks, nil
}
