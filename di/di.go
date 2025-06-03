package di

import (
	"context"
)

type (
	DiType string
	Scope  string
)

const (
	DiTypeSimpleDi = DiType("simpledi")
)

const (
	Singleton = Scope("singleton")
	Scoped    = Scope("scoped")
)

type DependencyFactoryFunc func(r Resolver) (any, error)

type Resolver interface {
	Get(key string) any
}

type Container interface {
	Resolver
	AddSingleton(key string, factory DependencyFactoryFunc)
	AddScoped(key string, factory DependencyFactoryFunc)
	Add(scope Scope, key string, factory DependencyFactoryFunc)
	WithScope(ctx context.Context) context.Context
}
