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

type valueInfo struct {
	value any
	once  sync.Once
}

var _ di.Container = (*container)(nil)

type container struct {
	dependencies map[string]dependencyInfo
	logger       logger.Logger
	mu           *sync.RWMutex
	root         *container
	parent       *container
	tracked      tracked
	values       map[string]*valueInfo
}

func NewSimpleDi(cfg *config.DiOptions) di.Container {
	c := &container{
		dependencies: map[string]dependencyInfo{},
		logger:       cfg.Logger,
		mu:           &sync.RWMutex{},
		tracked:      make(tracked),
		values:       map[string]*valueInfo{},
	}

	c.root = c

	return c
}

func (c *container) Add(scope di.Scope, key string, factory di.DependencyFactoryFunc) {
	utils.Assert(factory != nil, fmt.Errorf("di: the dependency with key `%s` can't be registered without a factory function", key))

	c.mu.Lock()
	defer c.mu.Unlock()

	c.dependencies[key] = dependencyInfo{
		key:     key,
		scope:   scope,
		factory: factory,
	}
	c.values[key] = &valueInfo{
		once: sync.Once{},
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
	c.logger.Debugf("di: trying to get the dependency with key: `%s`", key)

	info, exists := c.dependencies[key]
	if !exists {
		panic(fmt.Errorf("di: there is no dependency registered with the key: `%s`", key))
	}

	if _, exists := c.tracked[info.key]; exists {
		cycleChain := strings.Join(c.tracked.ordered(), " -> ")
		panic(fmt.Errorf("di: cycle dependencies building `%s`, chain: %s", info.key, cycleChain))
	}

	switch info.scope {
	case di.Singleton:
		return c.root.get(info)
	default:
		return c.get(info)
	}
}

func (c *container) get(info dependencyInfo) any {
	c.mu.RLock()
	valueInfo, exists := c.values[info.key]
	c.mu.RUnlock()

	utils.Assert(exists, fmt.Errorf("di: something went very wrong as the valueInfo for key `%s` is missing", info.key))

	var err error
	valueInfo.once.Do(func() {
		valueInfo.value, err = info.factory(c.builder(info))
	})

	if err != nil {
		c.mu.Lock()
		delete(c.values, info.key)
		c.mu.Unlock()
		panic(fmt.Errorf("di: failed to build the dependency `%s` with error: %w", info.key, err))
	}

	return valueInfo.value
}

func (c *container) scoped() *container {
	return &container{
		dependencies: c.dependencies,
		logger:       c.logger,
		mu:           &sync.RWMutex{},
		parent:       c,
		root:         c.root,
		tracked:      make(tracked),
		values:       map[string]*valueInfo{},
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
