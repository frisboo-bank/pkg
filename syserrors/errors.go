package syserrors

import "iter"

func MustBePositiveError(name string, value any) error {
	return Newf("%s must be >0: got %v", name, value)
}

func CantBeNegativeError(name string, value any) error {
	return Newf("%s cannot be negative: got %v", name, value)
}

func CantBeNilError(name string) error {
	return Newf("%s can't be nil", name)
}

func CantBeEmptyError(name string) error {
	return Newf("%s can't be empty", name)
}

func UnknownEnumError[T any](name string, enums iter.Seq[T]) error {
	return Newf("%s is invalid: available options %v", name, enums)
}
