package dig

import (
	"context"
	"fmt"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/container/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/waiter"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"
	waiterOptions "frisboo-bank/pkg/waiter/options"
	"reflect"
	"slices"
	"sync"

	"go.uber.org/dig"
)

const (
	hookGroupStartPattern = "hook-start-%s-%d"
	hookGroupStopPattern  = "hook-stop-%s-%d"
)

type digContainer struct {
	digContainer        *dig.Container
	digContainerOptions []dig.Option
	modules             []container.Module
	started             bool
	startOnce           sync.Once
	stopOnce            sync.Once
	waiter              waiterContracts.Waiter
	waiterOptions       []waiterOptions.WaiterOption
	metrics             *container.ContainerMetrics
	logger              loggerContracts.Logger
}

var _ contracts.Container = (*digContainer)(nil)

func NewDigContainer(modules []container.Module, logger loggerContracts.Logger) contracts.Container {
	return &digContainer{
		modules: modules,
		metrics: container.NewDigMetrics(),
		logger:  logger,
	}
}

func (d *digContainer) WithDigOptions(options ...dig.Option) *digContainer {
	d.digContainerOptions = options
	return d
}

func (d *digContainer) WithWaiterOptions(options ...waiterOptions.WaiterOption) contracts.Container {
	d.waiterOptions = options
	return d
}

func (d *digContainer) Start(ctx context.Context) error {
	var err error

	d.startOnce.Do(func() {
		err = d.start(ctx)
	})

	return err
}

func (d *digContainer) start(ctx context.Context) error {
	d.digContainer = dig.New(d.digContainerOptions...)

	waiterOpts := d.waiterOptions
	if waiterOpts == nil {
		waiterOpts = []waiterOptions.WaiterOption{
			waiterOptions.WithParentContext(ctx),
			waiterOptions.UseCancelOnShutdownSignal(),
		}
	}
	d.waiter = waiter.NewWaiter(waiterOpts...)

	modules := d.collectAllModules()

	for _, m := range modules {
		if err := d.registerProviders(m); err != nil {
			return fmt.Errorf("(dig-container) failed to register providers from module %s with error: %v", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerHooks(m); err != nil {
			return fmt.Errorf("(dig-container) failed to register hooks from module %s with error: %v", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerDecorators(m); err != nil {
			return fmt.Errorf("(dig-container) failed to register decorators from module %s with error: %v", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerInvokes(m); err != nil {
			return fmt.Errorf("(dig-container) failed to invoke functions from module %s with error: %v", m.GetName(), err)
		}
	}

	d.Logger().Debugf("(dig-container) module metrics:\n%s", d.metrics.ToString())

	for _, m := range modules {
		hooks, err := d.resolveHooks(m)
		if err != nil {
			return fmt.Errorf("(dig-container) failed to invoke hooks with error: %v", err)
		}

		d.waiter.Add(hooks...)
	}

	d.started = true
	return d.waiter.Wait()
}

func (d *digContainer) Stop(ctx context.Context) error {
	var err error

	d.stopOnce.Do(func() {
		err = d.stop(ctx)
	})

	return err
}

func (d *digContainer) stop(ctx context.Context) error {
	if !d.started {
		return fmt.Errorf("(dig-container) seams like there is no container running, you need to call Start first")
	}

	d.waiter.Cancel()
	return nil
}

func (d *digContainer) collectAllModules() []container.Module {
	modules := make([]container.Module, 0)
	queue := slices.Clone(d.modules)

	for len(queue) > 0 {
		module := queue[0]
		queue = queue[1:]

		modules = append(modules, module)

		if len(module.GetModules()) > 0 {
			d.metrics.IncrementModules(module.GetName(), len(module.GetModules()))

			queue = append(queue, module.GetModules()...)
		}
	}

	return modules
}

func (d *digContainer) registerProviders(module container.Module) error {
	d.metrics.IncrementProviders(module.GetName(), len(module.GetProviders()))

	for id, provider := range module.GetProviders() {
		options := filterOptions[dig.ProvideOption](provider.Options(), d.logger)

		if err := d.digContainer.Provide(provider.Fn(), options...); err != nil {
			return fmt.Errorf("failed to register provider %d with error: %v", id, err)
		}
	}

	return nil
}

func (d *digContainer) registerHooks(module container.Module) error {
	d.metrics.IncrementHooks(module.GetName(), len(module.GetHooks()))

	for id, hook := range module.GetHooks() {
		startGroup := fmt.Sprintf(hookGroupStartPattern, module.GetName(), id)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, module.GetName(), id)

		options := filterOptions[dig.ProvideOption](hook.Options(), d.logger)

		startOptions := append(options, dig.Group(startGroup))
		if err := d.digContainer.Provide(hook.StartFn(), startOptions...); err != nil {
			return fmt.Errorf("failed to register hook start %d with error: %v", id, err)
		}

		stopOptions := append(options, dig.Group(stopGroup))
		if err := d.digContainer.Provide(hook.StopFn(), stopOptions...); err != nil {
			return fmt.Errorf("failed to register hook stop %d with error: %v", id, err)
		}
	}

	return nil
}

func (d *digContainer) registerDecorators(module container.Module) error {
	d.metrics.IncrementDecorators(module.GetName(), len(module.GetDecorators()))

	for id, decorator := range module.GetDecorators() {
		options := filterOptions[dig.DecorateOption](decorator.Options(), d.logger)

		if err := d.digContainer.Decorate(decorator.Fn(), options...); err != nil {
			return fmt.Errorf("failed to register decorator %d with error: %v", id, err)
		}
	}

	return nil
}

func (d *digContainer) registerInvokes(module container.Module) error {
	d.metrics.IncrementInvokes(module.GetName(), len(module.GetInvokes()))

	for id, invoke := range module.GetInvokes() {
		options := filterOptions[dig.InvokeOption](invoke.Options(), d.logger)

		if err := d.digContainer.Invoke(invoke.Fn(), options...); err != nil {
			return fmt.Errorf("failed to invoke %d with error: %v", id, err)
		}
	}

	return nil
}

func (d *digContainer) resolveHooks(module container.Module) ([]waiterContracts.WaiterHook, error) {
	hooks := make([]waiterContracts.WaiterHook, 0)

	for i := range module.GetHooks() {
		startGroup := fmt.Sprintf(hookGroupStartPattern, module.GetName(), i)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, module.GetName(), i)

		hook := waiterContracts.WaiterHook{}

		startFns, err := resolveDynamicGroup[[]waiterContracts.WaitFunc](
			d.digContainer,
			startGroup,
			reflect.TypeOf([]waiterContracts.WaitFunc{}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve hook from group: %s", startGroup)
		}

		if len(startFns) > 0 {
			hook.Wait = startFns[0]
		}

		stopFns, err := resolveDynamicGroup[[]waiterContracts.CleanupFunc](
			d.digContainer,
			stopGroup,
			reflect.TypeOf([]waiterContracts.CleanupFunc{}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve hook from group: %s", stopGroup)
		}

		if len(stopFns) > 0 {
			hook.Cleanup = stopFns[0]
		}

		hooks = append(hooks, hook)
	}

	return hooks, nil
}

func (d *digContainer) Logger() loggerContracts.Logger {
	return d.logger
}
