package options

import (
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

type (
	DefaultFn[T any] func() T
	OptionFn[T any]  func(*T) error
)

// -----------------------------------------------------------------------------
// Option factories
// -----------------------------------------------------------------------------
func buildOption[T any, A any](fn func(*T, A) error) func(A) OptionFn[T] {
	return func(a A) OptionFn[T] {
		if fn == nil {
			return noopOption[T]()
		}
		return func(c *T) error { return fn(c, a) }
	}
}

func Option[T any, A any](fn func(*T, A)) func(A) OptionFn[T] {
	return buildOption(func(c *T, a A) error {
		fn(c, a)
		return nil
	})
}

func OptionErr[T any, A any](fn func(*T, A) error) func(A) OptionFn[T] {
	return buildOption(fn)
}

func buildVarOption[T any, A any](fn func(*T, ...A) error) func(...A) OptionFn[T] {
	return func(a ...A) OptionFn[T] {
		if fn == nil {
			return noopOption[T]()
		}
		return func(c *T) error { return fn(c, a...) }
	}
}

func VarOption[T any, A any](fn func(*T, ...A)) func(...A) OptionFn[T] {
	return buildVarOption(func(c *T, a ...A) error {
		fn(c, a...)
		return nil
	})
}

func VarOptionErr[T any, A any](fn func(*T, ...A) error) func(...A) OptionFn[T] {
	return buildVarOption(fn)
}

func noopOption[T any]() OptionFn[T] { return func(t *T) error { return nil } }

// -----------------------------------------------------------------------------
// Composition
// -----------------------------------------------------------------------------

func Apply[T any](target *T, opts ...OptionFn[T]) error {
	if target == nil {
		return syserrors.CantBeNilError("target")
	}

	for _, o := range opts {
		if err := o(target); err != nil {
			return err
		}
	}

	if v, ok := any(&target).(validation.Validatable); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func Compose[T any](opts ...OptionFn[T]) OptionFn[T] {
	if len(opts) == 0 {
		return noopOption[T]()
	}
	return func(t *T) error { return Apply(t, opts...) }
}
