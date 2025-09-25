package container

import (
	"context"
	"fmt"
	"sync"

	"frisboo-bank/pkg/container/contracts"
	"frisboo-bank/pkg/container/dependencies/module"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

var _ contracts.Container = (*container)(nil)

type container struct {
	adapter contracts.ContainerAdapter
	started bool
	mu      sync.Mutex
}

func New(adapter contracts.ContainerAdapter) contracts.Container {
	validation.Assert(adapter != nil, "container: your must set the adapter")

	return &container{
		adapter: adapter,
	}
}

func (c *container) RegisterModule(modules ...module.Module) error {
	modules, err := c.collectAllModules(modules...)
	if err != nil {
		return err
	}

	for _, module := range modules {
		for id, p := range module.Providers() {
			name := fmt.Sprintf("%s:provider:%d", module.Name(), id)
			if err := c.adapter.RegisterProvider(name, p); err != nil {
				return syserrors.Newf("provider %s registration failed with error: %w", name, err)
			}
		}
	}

	for _, module := range modules {
		for _, h := range module.Hooks() {
			if err := c.adapter.RegisterHook(h); err != nil {
				return syserrors.Newf("hook %s registration failed with error: %w", h.Name(), err)
			}
		}
	}

	for _, module := range modules {
		for id, d := range module.Decorators() {
			name := fmt.Sprintf("%s:decorator:%d", module.Name(), id)
			if err := c.adapter.RegisterDecorator(name, d); err != nil {
				return syserrors.Newf("decorator %s registration failed with error: %w", name, err)
			}
		}
	}

	for _, module := range modules {
		for id, i := range module.Invokers() {
			name := fmt.Sprintf("%s:invoker:%d", module.Name(), id)
			if err := c.adapter.RegisterInvoker(name, i); err != nil {
				return syserrors.Newf("invoker registration failed for module %q with error: %w", module.Name(), err)
			}
		}
	}

	return nil
}

func (c *container) Start(ctx context.Context) (err error) {
	validation.Assert(ctx != nil, "container: your must set the context")

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.started {
		return syserrors.Newf("container: already running")
	}

	if err := c.adapter.Start(ctx); err != nil {
		return err
	}

	c.started = true

	return nil
}

func (c *container) Stop(ctx context.Context) error {
	validation.Assert(ctx != nil, "container: your must set the context")

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.started {
		return syserrors.Newf("container: not running; call Start first")
	}

	err := c.adapter.Stop(ctx)
	c.started = false

	return err
}

func (c *container) Type() containertype.ContainerType {
	return c.adapter.Type()
}

func (c *container) collectAllModules(modules ...module.Module) ([]module.Module, error) {
	queue := modules
	tree := make([]module.Module, 0, len(queue))
	visited := make(map[module.Module]struct{})

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
				return nil, syserrors.Newf("module %q references itself; skipping to avoid cycle", module.Name())
			}
			queue = append(queue, child)
		}
	}

	return tree, nil
}
