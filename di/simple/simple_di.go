package simple

import (
	"fmt"
	"frisboo-bank/pkg/di"
	"sync"
)

type dependencyInfo struct {
	key     string
	scope   di.Scope
	factory di.DependencyFactoryFunc
}

type container struct {
	parent       *container
	dependencies map[string]dependencyInfo
	values       map[string]any
	tracked      any
	mu           sync.Mutex
}

func NewSimpleDi() di.Container {
	return &container{
		dependencies: make(map[string]dependencyInfo),
		values:       make(map[string]any),
	}
}

func (c *container) AddScoped(key string, dependencyFactoryFunc di.DependencyFactoryFunc) {
	c.Add(di.Scoped, key, dependencyFactoryFunc)
}

func (c *container) AddSingleton(key string, dependencyFactoryFunc di.DependencyFactoryFunc) {
	c.Add(di.Singleton, key, dependencyFactoryFunc)
}

func (c *container) Add(scope di.Scope, key string, dependencyFactoryFunc di.DependencyFactoryFunc) {
	c.dependencies[key] = dependencyInfo{
		key:     key,
		scope:   scope,
		factory: dependencyFactoryFunc,
	}
}

// // WithScope implements di.Container.
// func (c *container) WithScope(ctx context.Context) context.Context {
// 	return context.WithValue(ctx, c, val any)
// }

func (c *container) Get(key string) (any, error) {
	info, exists := c.dependencies[key]
	if !exists {
		return nil, fmt.Errorf("Di: there is dependency register with the key: %s", key)
	}

	if info.scope == di.Singleton {
		return c.getFromParentContainer(info)
	}

	return c.get(info)
}

func (c *container) getFromParentContainer(info dependencyInfo) (any, error) {
	if c.parent != nil {
		return c.parent.getFromParentContainer(info)
	}

	return c.get(info)
}

func (c *container) get(info dependencyInfo) (any, error) {
	c.mu.Lock()

	value, exists := c.values[info.key]
	if !exists {
	}
}
