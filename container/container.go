package container

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/utils"
)

var _ contracts.Container = (*container)(nil)

type container struct {
	cfg       *config.Config
	adapter   contracts.ContainerAdapter
	logger    loggerContracts.Logger
	modules   []contracts.Module
	started   bool
	startOnce sync.Once
	stopOnce  sync.Once
}

func New(
	adapter contracts.ContainerAdapter,
	logger loggerContracts.Logger,
	opts *options.OptionBuilder[config.Config],
) (contracts.Container, error) {
	utils.Assert(adapter != nil, "container: your must set the adapter")
	utils.Assert(logger != nil, "container: your must set the logger")
	utils.Assert(opts != nil, "container: your must set the opts")

	cfg := opts.Build()

	container := &container{
		cfg:     cfg,
		adapter: adapter,
		logger:  logger,
	}

	if err := adapter.Setup(cfg); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *container) RegisterModule(modules ...contracts.Module) error {
	modules, err := c.collectAllModules(modules...)
	if err != nil {
		return err
	}

	for _, module := range modules {
		if err := c.adapter.RegisterProvider(module.Providers()...); err != nil {
			return fmt.Errorf("provider registration failed for module %q with error: %w", module.Name(), err)
		}
	}

	for _, module := range modules {
		if err := c.adapter.RegisterHook(module.Hooks()...); err != nil {
			return fmt.Errorf("hook registration failed for module %q with error: %w", module.Name(), err)
		}
	}

	for _, module := range modules {
		if err := c.adapter.RegisterDecorator(module.Decorators()...); err != nil {
			return fmt.Errorf("decorator registration failed for module %q with error: %w", module.Name(), err)
		}
	}

	for _, module := range modules {
		if err := c.adapter.RegisterInvoker(module.Invokers()...); err != nil {
			return fmt.Errorf("invoker registration failed for module %q with error: %w", module.Name(), err)
		}
	}

	return nil
}

func (c *container) Start(ctx context.Context) (err error) {
	utils.Assert(ctx != nil, "container: your must set the context")
	utils.Assert(!c.started, "container: container already running")

	c.startOnce.Do(func() {
		err = errors.Join(err, c.adapter.Start(ctx))
		if err == nil {
			c.started = true
		}
	})

	return
}

func (c *container) Stop(ctx context.Context) (err error) {
	utils.Assert(ctx != nil, "container: your must set the context")
	utils.Assert(c.started, "container: no container running; you must call Start first")

	c.stopOnce.Do(func() {
		err = errors.Join(err, c.adapter.Stop(ctx))
		c.started = false
	})

	return
}

func (c *container) Type() containertype.ContainerType {
	return c.adapter.Type()
}

func (c *container) collectAllModules(modules ...contracts.Module) ([]contracts.Module, error) {
	queue := modules
	tree := make([]contracts.Module, 0, len(queue))
	visited := make(map[contracts.Module]struct{})

	for len(queue) > 0 {
		module := queue[0]
		queue = queue[1:]

		// Check for cycles/duplicates based on pointer identity.
		if _, seen := visited[module]; seen {
			continue
		}
		visited[module] = struct{}{}

		tree = append(tree, module)

		for _, child := range module.Modules() {
			// Defensive: avoid self-cycle.
			if child == module {
				return nil, fmt.Errorf("module %q references itself; skipping to avoid cycle", module.Name())
			}
			queue = append(queue, child)
		}
	}

	return tree, nil
}
