package registry

import (
	configLoaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/syserrors"

	"dario.cat/mergo"
)

type Registry[T any] struct {
	key      string
	kind     string
	baseline func() T

	Default   T            `mapstructure:",squash"`
	Instances map[string]T `mapstructure:"instances"`
}

func Load[T any](
	configLoader configLoaderContracts.ConfigLoader,
	env environment.Environment,
	key string,
	kind string,
	baseline func() T,
) (*Registry[T], error) {
	reg := &Registry[T]{
		key:      key,
		kind:     kind,
		baseline: baseline,
	}

	if err := configLoader.LoadKey(env, reg, key); err != nil {
		return nil, err
	}

	return reg, nil
}

func (r *Registry[T]) Has(name string) bool {
	if len(r.Instances) == 0 {
		return false
	}
	_, ok := r.Instances[name]
	return ok
}

func (r *Registry[T]) GetByNameOrDefault(name string) (T, error) {
	if name != "" {
		return r.GetByName(name)
	}
	return r.GetDefault()
}

func (r *Registry[T]) GetDefault() (T, error) {
	var zero T

	base := r.baseline()
	if err := mergo.Merge(&base, r.Default); err != nil {
		return zero, syserrors.Wrapf(err, "merge default")
	}

	return base, nil
}

func (r *Registry[T]) GetByName(name string) (T, error) {
	var zero T
	if name == "" {
		return zero, syserrors.CantBeEmptyError("name")
	}

	inst, ok := r.Instances[name]
	if !ok {
		return zero, syserrors.Newf("%s %q not found in the config", r.kind, name)
	}

	base, err := r.GetDefault()
	if err != nil {
		return zero, err
	}

	if err := mergo.Merge(&base, inst); err != nil {
		return zero, syserrors.Wrapf(err, "merge instance %q", name)
	}

	return base, nil
}
