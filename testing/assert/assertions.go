package assert

import (
	"fmt"
	"unicode"

	"github.com/stretchr/testify/assert"
)

type tHelper = interface {
	Helper()
}

func IsAllDigit(t assert.TestingT, object string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(tHelper); ok {
		h.Helper()
	}

	value := true
	for _, r := range object {
		if !unicode.IsDigit(r) {
			value = false
			break
		}
	}

	if value {
		return true
	}

	return assert.Fail(t, fmt.Sprintf("Expected all digits, but got: %#v", object), msgAndArgs...)
}
