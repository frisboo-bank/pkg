package validation

import (
	"frisboo-bank/pkg/reflection"
	"frisboo-bank/pkg/syserrors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// EnumOneOf validates that the value is not invalid or unknown
func EnumOneOf[T reflection.EnumValue](options reflection.EnumContainer[T]) validation.RuleFunc {
	var opts []T
	for o := range options.All() {
		opts = append(opts, o)
	}

	return func(value any) error {
		v, ok := (value).(T)
		if !ok {
			return syserrors.New("Must be a valid enum")
		}

		if !v.IsValid() {
			return syserrors.MustBeOneOf("", v, opts)
		}

		return nil
	}
}
