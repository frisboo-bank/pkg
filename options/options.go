package options

import (
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

type (
	// DefaultFn returns a default-initialized T (e.g., with sensible defaults).
	DefaultFn[T any] func() T
	// OptionFn mutates the target in place and can fail.
	OptionFn[T any] func(*T) error
)

// -----------------------------------------------------------------------------
// Option factories
// -----------------------------------------------------------------------------

// buildOption creates a typed option factory for a single-arg setter.
func buildOption[T any, A any](fn func(*T, A) error) func(A) OptionFn[T] {
	return func(a A) OptionFn[T] {
		if fn == nil {
			return noopOption[T]()
		}
		return func(c *T) error { return fn(c, a) }
	}
}

// Option wraps a setter that cannot fail.
func Option[T any, A any](fn func(*T, A)) func(A) OptionFn[T] {
	return buildOption(func(c *T, a A) error {
		fn(c, a)
		return nil
	})
}

// OptionErr wraps a setter that can fail.
func OptionErr[T any, A any](fn func(*T, A) error) func(A) OptionFn[T] {
	return buildOption(fn)
}

// buildVarOption creates a typed option factory for a varargs setter.
func buildVarOption[T any, A any](fn func(*T, ...A) error) func(...A) OptionFn[T] {
	return func(a ...A) OptionFn[T] {
		if fn == nil {
			return noopOption[T]()
		}
		return func(c *T) error { return fn(c, a...) }
	}
}

// VarOption wraps a varargs setter that cannot fail.
func VarOption[T any, A any](fn func(*T, ...A)) func(...A) OptionFn[T] {
	return buildVarOption(func(c *T, a ...A) error {
		fn(c, a...)
		return nil
	})
}

// VarOptionErr wraps a varargs setter that can fail.
func VarOptionErr[T any, A any](fn func(*T, ...A) error) func(...A) OptionFn[T] {
	return buildVarOption(fn)
}

// OptionIf applies opt only when cond is true; otherwise it no-ops.
func OptionIf[T any](cond bool, opt OptionFn[T]) OptionFn[T] {
	if cond {
		return opt
	}
	return noopOption[T]()
}

// noopOption returns an option that does nothing and succeeds.
func noopOption[T any]() OptionFn[T] { return func(*T) error { return nil } }

// -----------------------------------------------------------------------------
// Composition
// -----------------------------------------------------------------------------

// Apply applies opts in order and then validates target if it implements
func Apply[T any](target *T, opts ...OptionFn[T]) error {
	if target == nil {
		return syserrors.CantBeNilError("target")
	}

	for _, o := range opts {
		if err := o(target); err != nil {
			return err
		}
	}

	if v, ok := any(target).(validation.Validatable); ok {
		return v.Validate()
	}

	return nil
}

func Compose[T any](opts ...OptionFn[T]) OptionFn[T] {
	if len(opts) == 0 {
		return noopOption[T]()
	}
	return func(t *T) error { return Apply(t, opts...) }
}
