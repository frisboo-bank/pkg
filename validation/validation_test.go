package validation

import (
	"errors"
	"slices"
	"strings"
	"testing"

	enum "frisboo-bank/pkg/validation/enums"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestEnsureEven(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		customEvenError := errors.New("custom error")
		v := rapid.Int().Draw(t, "v")

		err := Ensure(
			Field[int]{Name: "evenField", Value: v},
			func(v int) bool { return v%2 == 0 },
			func(Field[int]) error { return customEvenError },
		)
		if v%2 == 0 {
			assert.NoError(t, err, "even number %d should pass", v)
		} else {
			assert.ErrorIs(t, err, customEvenError, "odd number %d should fail", v)
		}
	})
}

func TestCustomErr(t *testing.T) {
	customError := errors.New("custom error")

	err := Positive("v", -1, func(f Field[int]) error {
		assert.Equal(t, "v", f.Name)
		assert.Equal(t, -1, f.Value)

		return customError
	})

	assert.ErrorIs(t, err, customError)
}

func TestNotEmpty(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		v := rapid.String().Draw(rt, "v")

		err := NotEmpty("v", v)

		if strings.TrimSpace(v) == "" {
			assert.Error(rt, err, "expected error for empty/whitespace string %q", v)
		} else {
			assert.NoError(rt, err, "expected no error for non-empty string %q", v)
		}
	})
}

func TestPositive(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		v := rapid.Int().Draw(rt, "v")

		err := Positive("v", v)

		if v > 0 {
			assert.NoError(rt, err, "v > 0 should pass")
		} else {
			assert.Error(rt, err, "v <= 0 should error")
		}
	})
}

func TestNonNegative(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		v := rapid.Int().Draw(rt, "v")

		err := NonNegative("v", v)

		if v >= 0 {
			assert.NoError(rt, err, "v >= 0 should pass")
		} else {
			assert.Error(rt, err, "v < 0 should error")
		}
	})
}

func TestIsTrue(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		v := rapid.Bool().Draw(rt, "v")

		err := IsTrue("v", v)

		if v {
			assert.NoError(rt, err, "true should pass")
		} else {
			assert.Error(rt, err, "false should error")
		}
	})
}

func TestIsFalse(t *testing.T) {
	rapid.Check(t, func(rt *rapid.T) {
		v := rapid.Bool().Draw(rt, "v")

		err := IsFalse("v", v)

		if !v {
			assert.NoError(rt, err, "false should pass")
		} else {
			assert.Error(rt, err, "true should error")
		}
	})
}

func TestOneOf(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		options := rapid.SliceOfNDistinct(rapid.Int(), 0, 15, rapid.ID[int]).
			Draw(t, "options")

		wantMember := rapid.Bool().Draw(t, "wantMember")
		if len(options) == 0 {
			wantMember = false
		}

		var v int
		if wantMember {
			idx := rapid.IntRange(0, len(options)-1).Draw(t, "idx")
			v = options[idx]
		} else {
			for {
				c := rapid.Int().Draw(t, "candidate")
				if !slices.Contains(options, c) {
					v = c
					break
				}
			}
		}

		err := OneOf("v", v, options)

		if wantMember {
			assert.NoError(t, err, "expected %d to be in %v", v, options)
		} else {
			assert.Error(t, err, "expected %d NOT to be in %v", v, options)
		}
	})
}

func TestEnumOneOf(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		var options []enums.EnumType
		for o := range enum.EnumTypes.All() {
			options = append(options, o)
		}

		invalid := rapid.Bool().Draw(t, "invalid")
		unknown := rapid.Bool().Draw(t, "unknown")

		var v enum.EnumType
		if invalid || unknown {
			v = enum.EnumTypes.UNKNOWN
		} else {
			idx := rapid.IntRange(0, len(options)-1).
				Draw(t, "idx")
			v = options[idx]
		}

		err := EnumOneOf("v", v, enum.EnumTypes)

		if invalid || unknown {
			assert.Error(t, err, "expected %d NOT to be in %v", v, options)
		} else {
			assert.NoError(t, err, "expected %d to be in %v", v, options)
		}
	})
}
