package simpledi

import (
	"context"
	"fmt"
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/di/config"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/utils"
	"strings"
	"sync"
)

type contextKey int

const containerKey contextKey = 1

type dependencyInfo struct {
	factory di.DependencyFactoryFunc
	key     string
	scope   di.Scope
}

var _ di.Container = (*container)(nil)

type container struct {
	dependencies map[string]dependencyInfo
	logger       logger.Logger
	mu           *sync.Mutex
	root         *container
	parent       *container
	tracked      tracked
	values       map[string]any
}

type buildingDependencyChan = chan struct{}

func NewSimpleDi(cfg *config.DiOptions) di.Container {
	c := &container{
		dependencies: map[string]dependencyInfo{},
		logger:       cfg.Logger,
		values:       map[string]any{},
	}

	c.root = c

	return c
}

func (c *container) Add(scope di.Scope, key string, factory di.DependencyFactoryFunc) {
	utils.Assert(factory == nil, fmt.Errorf("di: the dependency with key %s can't be registered without a factory function", key))
	c.mu.Lock()
	defer c.mu.Unlock()

	c.dependencies[key] = dependencyInfo{
		key:     key,
		scope:   scope,
		factory: factory,
	}
}

func (c *container) AddScoped(key string, factory di.DependencyFactoryFunc) {
	c.Add(di.Scoped, key, factory)
}

func (c *container) AddSingleton(key string, factory di.DependencyFactoryFunc) {
	c.Add(di.Singleton, key, factory)
}

func (c *container) WithScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, containerKey, c.scoped())
}

func (c *container) Get(key string) any {
	c.logger.Debugf("di: trying to get the dependency with key: %s", key)

	info, exists := c.dependencies[key]
	if !exists {
		panic(fmt.Errorf("di: there is no dependency registered with the key: %s", key))
	}

	if _, exists := c.tracked[info.key]; exists {
		cycleChain := strings.Join(c.tracked.ordered(), " -> ")
		panic(fmt.Errorf("di: cycle dependencies building %s, chain: %s", info.key, cycleChain))
	}

	switch info.scope {
	case di.Singleton:
		return c.root.get(info)
	default:
		return c.get(info)
	}
}

func (c *container) get(info dependencyInfo) any {
	c.mu.Lock()

	value, exists := c.values[info.key]
	if !exists {
		c.logger.Debugf("di: the dependency %s is not yet created and will be now created", info.key)

		tempValue := make(buildingDependencyChan)
		c.values[info.key] = tempValue
		c.mu.Unlock()
		return c.build(info, tempValue)
	}

	c.mu.Unlock()

	if tempValue, isTemp := value.(buildingDependencyChan); isTemp {
		c.logger.Debugf("di: the dependency %s is pending for creation and will be returned when done", info.key)

		<-tempValue

		return c.get(info)
	}

	c.logger.Debugf("di: the dependency %s is already created and will be returned", info.key)
	return value
}

func (c *container) build(info dependencyInfo, tempValue buildingDependencyChan) any {
	builder := c.builder(info)

	value, err := info.factory(builder)
	if err != nil {
		c.mu.Lock()
		delete(c.values, info.key)
		c.mu.Unlock()
		close(tempValue)
		panic(fmt.Errorf("di: failed to buid the dependency %s with error: %w", info.key, err))
	}

	c.mu.Lock()
	c.values[info.key] = value
	c.mu.Unlock()
	close(tempValue)

	return value
}

func (c *container) scoped() *container {
	return &container{
		dependencies: c.dependencies,
		logger:       c.logger,
		parent:       c,
		root:         c.root,
		values:       make(map[string]any),
		tracked:      make(tracked),
	}
}

func (c *container) builder(info dependencyInfo) *container {
	return &container{
		dependencies: c.dependencies,
		logger:       c.logger,
		mu:           c.mu,
		root:         c.root,
		parent:       c.parent,
		tracked:      c.tracked.add(info),
		values:       c.values,
	}
}
