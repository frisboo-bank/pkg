package di

type Scope string

const (
	Singleton = Scope("singleton")
	Scoped    = Scope("scoped")
)

type DependencyFactoryFunc func(c Container) (any, error)

type Container interface {
	AddSingleton(key string, fn DependencyFactoryFunc)
	AddScoped(key string, fn DependencyFactoryFunc)
	Add(scope Scope, key string, fn DependencyFactoryFunc)
	// WithScope(ctx context.Context) context.Context
	Get(key string) (any, error)
}
