package validation

import (
	"testing"

	testenum "frisboo-bank/pkg/validation/testdata/enums/test_enum"

	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestEnumOneOf(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		var options []testenum.TestEnum
		for o := range testenum.TestEnums.All() {
			options = append(options, o)
		}

		invalid := rapid.Bool().Draw(t, "invalid")

		var v testenum.TestEnum
		switch {
		case invalid:
			v = testenum.TestEnums.UNKNOWN
		default:
			idx := rapid.IntRange(0, len(options)-1).Draw(t, "idx")
			v = options[idx]
		}

		err := EnumOneOf(testenum.TestEnums)(v)

		if invalid {
			assert.Error(t, err, "expected %d NOT to be in %v", v, options)
		} else {
			assert.NoError(t, err, "expected %d to be in %v", v, options)
		}
	})
}
