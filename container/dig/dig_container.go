package dig

import (
	"context"
	"fmt"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/container/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/utils"
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

// NewDigContainer creates a new digContainer with the given modules and logger.
// The container is not started until Start is called.
func NewDigContainer(modules []container.Module, logger loggerContracts.Logger) contracts.Container {
	utils.Assert(logger != nil, "(dig-container) logger must not be nil")

	return &digContainer{
		modules: modules,
		metrics: container.NewDigMetrics(),
		logger:  logger,
	}
}

// WithDigOptions sets dig.Container options for the container.
func (d *digContainer) WithDigOptions(options ...dig.Option) *digContainer {
	d.digContainerOptions = options
	return d
}

// WithWaiterOptions sets Waiter options for the container.
func (d *digContainer) WithWaiterOptions(options ...waiterOptions.WaiterOption) contracts.Container {
	d.waiterOptions = options
	return d
}

// Start initializes and starts all modules, providers, hooks, decorators, and invokes.
// It is safe to call Start multiple times; initialization will only happen once.
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
			return fmt.Errorf("(dig-container) failed to register providers from module %s: %w", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerHooks(m); err != nil {
			return fmt.Errorf("(dig-container) failed to register hooks from module %s: %w", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerDecorators(m); err != nil {
			return fmt.Errorf("(dig-container) failed to register decorators from module %s: %w", m.GetName(), err)
		}
	}

	for _, m := range modules {
		if err := d.registerInvokes(m); err != nil {
			return fmt.Errorf("(dig-container) failed to invoke functions from module %s: %w", m.GetName(), err)
		}
	}

	d.Logger().Debugf("(dig-container) module metrics:\n%s", d.metrics.ToString())

	for _, m := range modules {
		hooks, err := d.resolveHooks(m)
		if err != nil {
			return fmt.Errorf("(dig-container) failed to invoke hooks from module %s: %w", m.GetName(), err)
		}

		d.waiter.Add(hooks...)
	}

	d.started = true
	return d.waiter.Wait()
}

// Stop gracefully stops the running container and its modules.
// It is safe to call Stop multiple times; teardown will only happen once.
func (d *digContainer) Stop(ctx context.Context) error {
	var err error
	d.stopOnce.Do(func() {
		err = d.stop(ctx)
	})
	return err
}

func (d *digContainer) stop(ctx context.Context) error {
	if !d.started {
		return fmt.Errorf("(dig-container) no container running; you must call Start first")
	}
	d.waiter.Cancel()
	return nil
}

// collectAllModules traverses all modules using BFS, detects cycles, and avoids duplicate registrations.
// It uses pointer identity to ensure each module is only processed once.
func (d *digContainer) collectAllModules() []container.Module {
	modules := make([]container.Module, 0)
	queue := slices.Clone(d.modules)
	visited := make(map[container.Module]struct{})

	for len(queue) > 0 {
		module := queue[0]
		queue = queue[1:]

		// Check for cycles/duplicates based on pointer identity.
		if _, seen := visited[module]; seen {
			continue
		}
		visited[module] = struct{}{}

		modules = append(modules, module)

		for _, child := range module.GetModules() {
			// Defensive: avoid self-cycle.
			if child == module {
				d.Logger().Warnf("Module %q references itself; skipping to avoid cycle.", module.GetName())
				continue
			}
			queue = append(queue, child)
		}
	}

	return modules
}

func (d *digContainer) registerProviders(module container.Module) error {
	d.metrics.IncrementProviders(module.GetName(), len(module.GetProviders()))

	for id, provider := range module.GetProviders() {
		options, err := filterOptions[dig.ProvideOption](provider.Options())
		if err != nil {
			return fmt.Errorf("failed to convert options for provider %d in module %s: %w", id, module.GetName(), err)
		}

		if err := d.digContainer.Provide(provider.Fn(), options...); err != nil {
			return fmt.Errorf("failed to register provider %d in module %s: %w", id, module.GetName(), err)
		}
	}

	return nil
}

func (d *digContainer) registerHooks(module container.Module) error {
	d.metrics.IncrementHooks(module.GetName(), len(module.GetHooks()))

	for id, hook := range module.GetHooks() {
		startGroup := fmt.Sprintf(hookGroupStartPattern, module.GetName(), id)
		stopGroup := fmt.Sprintf(hookGroupStopPattern, module.GetName(), id)

		options, err := filterOptions[dig.ProvideOption](hook.Options())
		if err != nil {
			return fmt.Errorf("failed to convert options for hook %d in module %s: %w", id, module.GetName(), err)
		}

		startOptions := append(options, dig.Group(startGroup))
		if err := d.digContainer.Provide(hook.StartFn(), startOptions...); err != nil {
			return fmt.Errorf("failed to register hook start %d in module %s: %w", id, module.GetName(), err)
		}

		stopOptions := append(options, dig.Group(stopGroup))
		if err := d.digContainer.Provide(hook.StopFn(), stopOptions...); err != nil {
			return fmt.Errorf("failed to register hook stop %d in module %s: %w", id, module.GetName(), err)
		}
	}

	return nil
}

func (d *digContainer) registerDecorators(module container.Module) error {
	d.metrics.IncrementDecorators(module.GetName(), len(module.GetDecorators()))

	for id, decorator := range module.GetDecorators() {
		options, err := filterOptions[dig.DecorateOption](decorator.Options())
		if err != nil {
			return fmt.Errorf("failed to convert options for decorator %d in module %s: %w", id, module.GetName(), err)
		}

		if err := d.digContainer.Decorate(decorator.Fn(), options...); err != nil {
			return fmt.Errorf("failed to register decorator %d in module %s: %w", id, module.GetName(), err)
		}
	}

	return nil
}

func (d *digContainer) registerInvokes(module container.Module) error {
	d.metrics.IncrementInvokes(module.GetName(), len(module.GetInvokes()))

	for id, invoke := range module.GetInvokes() {
		options, err := filterOptions[dig.InvokeOption](invoke.Options())
		if err != nil {
			return fmt.Errorf("failed to convert options for invoke %d in module %s: %w", id, module.GetName(), err)
		}

		if err := d.digContainer.Invoke(invoke.Fn(), options...); err != nil {
			return fmt.Errorf("failed to invoke %d in module %s: %w", id, module.GetName(), err)
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
			return nil, fmt.Errorf("failed to retrieve start hook from group %s in module %s: %w", startGroup, module.GetName(), err)
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
			return nil, fmt.Errorf("failed to retrieve stop hook from group %s in module %s: %w", stopGroup, module.GetName(), err)
		}

		if len(stopFns) > 0 {
			hook.Cleanup = stopFns[0]
		}

		hooks = append(hooks, hook)
	}

	return hooks, nil
}

// Logger returns the logger for this container.
func (d *digContainer) Logger() loggerContracts.Logger {
	return d.logger
}
