package validation

import (
	"slices"
	"strings"

	"frisboo-bank/pkg/reflection"
	"frisboo-bank/pkg/syserrors"
)

// Field represents a named value under validation.
type Field[T any] struct {
	Name  string
	Value T
}

// ErrorFn builds an error for a field which failed validation.
type ErrorFn[T any] = func(f Field[T]) error

// Ensure applies predicate to the field value; if predicate returns true it returns nil,
// otherwise it calls errFn(field).
func Ensure[T any](f Field[T], predicate func(v T) bool, errFn ErrorFn[T]) error {
	if predicate(f.Value) {
		return nil
	}
	return errFn(f)
}

// NotNil validates that value is not nil (including typed-nil interface cases).
func NotNil[T any](name string, value T, errFn ...ErrorFn[T]) error {
	if !reflection.IsPossiblyNil(value) {
		return nil
	}

	return Ensure[T](
		Field[T]{Name: name, Value: value},
		func(v T) bool { return !reflection.IsNil(any(v)) },
		firstErr(errFn, func(f Field[T]) error {
			return syserrors.CantBeNilError(f.Name)
		}),
	)
}

// NotEmpty validates that a string (after TrimSpace) is non-empty.
func NotEmpty(name string, value string, errFn ...ErrorFn[string]) error {
	return Ensure[string](
		Field[string]{Name: name, Value: value},
		func(v string) bool { return strings.TrimSpace(v) != "" },
		firstErr(errFn, func(f Field[string]) error {
			return syserrors.CantBeEmptyError(f.Name)
		}),
	)
}

// Positive validates that a numeric value is > 0.
func Positive[T reflection.Numeric](name string, value T, errFn ...ErrorFn[T]) error {
	return Ensure[T](
		Field[T]{Name: name, Value: value},
		func(v T) bool { return v > 0 },
		firstErr(errFn, func(f Field[T]) error {
			return syserrors.MustBePositiveError(f.Name, f.Value)
		}),
	)
}

// NonNegative validates that a numeric value is >= 0.
func NonNegative[T reflection.Numeric](name string, value T, errFn ...ErrorFn[T]) error {
	return Ensure[T](
		Field[T]{Name: name, Value: value},
		func(v T) bool { return v >= 0 },
		firstErr(errFn, func(f Field[T]) error {
			return syserrors.CantBeNegativeError(f.Name, f.Value)
		}),
	)
}

// IsTrue validates that a boolean is true.
func IsTrue(name string, value bool, errFn ...ErrorFn[bool]) error {
	return Ensure[bool](
		Field[bool]{Name: name, Value: value},
		func(v bool) bool { return v },
		firstErr(errFn, func(f Field[bool]) error {
			return syserrors.MustBeTrue(f.Name)
		}),
	)
}

// IsFalse validates that a boolean is false.
func IsFalse(name string, value bool, errFn ...ErrorFn[bool]) error {
	return Ensure[bool](
		Field[bool]{Name: name, Value: value},
		func(v bool) bool { return !v },
		firstErr(errFn, func(f Field[bool]) error {
			return syserrors.MustBeFalse(f.Name)
		}),
	)
}

func OneOf[T comparable](name string, value T, options []T, errFn ...ErrorFn[T]) error {
	return Ensure[T](
		Field[T]{Name: name, Value: value},
		func(v T) bool { return slices.Contains(options, v) },
		firstErr(errFn, func(f Field[T]) error {
			return syserrors.MustBeOneOf(f.Name, f.Value, options)
		}),
	)
}

func EnumOneOf[T reflection.EnumValue](name string, value T, options reflection.EnumContainer[T], errFn ...ErrorFn[T]) error {
	return Ensure[T](
		Field[T]{Name: name, Value: value},
		func(v T) bool { return value.IsValid() },
		firstErr(errFn, func(f Field[T]) error {
			return syserrors.MustBeOneOf(f.Name, f.Value, nil)
		}),
	)
}

// firstErr returns the first provided custom error function if present, otherwise fallback.
func firstErr[T any](candidates []ErrorFn[T], fallback ErrorFn[T]) ErrorFn[T] {
	if len(candidates) > 0 && candidates[0] != nil {
		return candidates[0]
	}
	return fallback
}
