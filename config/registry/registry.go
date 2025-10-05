package registry

import (
	"reflect"
	"strings"

	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	configLoaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	responseformat "frisboo-bank/pkg/health/enums/response_format"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	encodingtype "frisboo-bank/pkg/logger/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	migrationcommandtype "frisboo-bank/pkg/migration/enums/migration_command_type"
	migratortype "frisboo-bank/pkg/migration/enums/migrator_type"
	"frisboo-bank/pkg/options"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
	"frisboo-bank/pkg/syserrors"

	"dario.cat/mergo"
)

var _ Registry[any] = (*registry[any])(nil)

func NameNotFoundError(kind string, name string, examples []string) error {
	err := syserrors.Newf("kind:{%s} name:{%s} not found in the config", kind, name)

	if len(examples) == 0 {
		return syserrors.Wrapf(err, "no kind:{%s} registered", kind)
	}
	return syserrors.Wrapf(err, "wrong name:{%s} only:{%s} are available", name, strings.Join(examples, ", "))
}

type Registry[T any] interface {
	Has(name string) bool
	Names() []string
	GetDefault(opts ...options.OptionFn[T]) (T, error)
	GetByName(name string, opts ...options.OptionFn[T]) (T, error)
}

type registry[T any] struct {
	key      string
	kind     string
	baseline func() T

	DefaultConfig T            `mapstructure:",squash"`
	Instances     map[string]T `mapstructure:"instances"`
}

func Load[T any](
	configLoader configLoaderContracts.ConfigLoader,
	env environment.Environment,
	key string,
	kind string,
	baseline func() T,
) (registry[T], error) {
	var zero registry[T]

	reg := registry[T]{
		key:      key,
		kind:     kind,
		baseline: baseline,
	}

	if err := configLoader.LoadKey(env, &reg, key); err != nil {
		return zero, err
	}

	return reg, nil
}

func (r registry[T]) Has(name string) bool {
	if len(r.Instances) == 0 {
		return false
	}
	_, ok := r.Instances[name]
	return ok
}

func (r registry[T]) GetDefault(opts ...options.OptionFn[T]) (T, error) {
	var zero T

	base := r.baseline()

	if err := mergeConfig(&base, r.DefaultConfig); err != nil {
		return zero, syserrors.Wrapf(err, "merge default")
	}

	if err := options.Apply(&base, opts...); err != nil {
		return zero, syserrors.Wrapf(err, "apply options")
	}

	return base, nil
}

func (r registry[T]) GetByName(name string, opts ...options.OptionFn[T]) (T, error) {
	var zero T

	if name == "" {
		return zero, syserrors.CantBeEmptyError("name")
	}

	inst, ok := r.Instances[name]
	if !ok {
		return zero, NameNotFoundError(r.kind, name, r.Names())
	}

	base := r.baseline()

	if err := mergeConfig(&base, r.DefaultConfig); err != nil {
		return zero, syserrors.Wrapf(err, "merge default %s", name)
	}

	if err := mergeConfig(&base, inst); err != nil {
		return zero, syserrors.Wrapf(err, "merge instance %s", name)
	}

	if err := options.Apply(&base, opts...); err != nil {
		return zero, syserrors.Wrapf(err, "apply options %s", name)
	}

	return base, nil
}

func (r registry[T]) Names() []string {
	ns := make([]string, 0, len(r.Instances))
	for name := range r.Instances {
		ns = append(ns, name)
	}
	return ns
}

type enumTransformer struct{}

type enum interface {
	IsValid() bool
}

// Transformer implements mergo.Transformers.
func (e enumTransformer) Transformer(t reflect.Type) func(dst reflect.Value, src reflect.Value) error {
	isZero := func(v reflect.Value) bool {
		return v.IsZero()
	}

	switch t {
	case
		reflect.TypeOf(cachetype.CacheType{}),
		reflect.TypeOf(encodingtype.EncodingType{}),
		reflect.TypeOf(httpservertype.HttpServerType{}),
		reflect.TypeOf(responseformat.ResponseFormat{}),
		reflect.TypeOf(loggertype.LoggerType{}),
		reflect.TypeOf(loglevel.LogLevel{}),
		reflect.TypeOf(migratortype.MigratorType{}),
		reflect.TypeOf(migrationcommandtype.MigrationCommandType{}),
		reflect.TypeOf(rpcservertype.RpcServerType{}):
		return func(dst, src reflect.Value) error {
			if isZero(src) {
				return nil
			}
			// Otherwise, same semantics as WithOverride
			dst.Set(src)
			return nil
		}
	}

	return nil
}

func mergeConfig[T any](dst *T, src T) error {
	return mergo.Merge(dst, src,
		mergo.WithOverride,
		mergo.WithTypeCheck,
		mergo.WithTransformers(enumTransformer{}),
	)
}
