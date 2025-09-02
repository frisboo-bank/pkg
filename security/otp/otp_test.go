package otp

import (
	"testing"

	cAssert "frisboo-bank/pkg/testing/assert"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWithLength(t *testing.T) {
	for _, length := range []int{1, 4, 6, 8} {
		code, err := GenerateWithLength(length)

		assert.Nil(t, err, "GenerateWithLength(%d) unexpected error: %v", length, err)
		assert.Len(t, code, length, "GenerateWithLength(%d) wrong length: got %d", length, len(code))
		cAssert.IsAllDigit(t, code, "GenerateWithLength(%d) found non-digit: got %q", length, code)
	}
}

func TestGenerate(t *testing.T) {
	code, err := Generate()

	assert.Nil(t, err, "Generate() unexpected error: %w", err)
	assert.Len(t, code, DefaultLength, "Generate() wrong length: got %d", len(code))
}

func TestErrorLengthOutOfRange(t *testing.T) {
	code, err := GenerateWithLength(0)

	assert.Equal(t, "", code, "GenerateWithLength(%d) unexpected code generated: got %q", 0, code)
	assert.EqualError(t, err, "got 0: length out of allowed range")
}

// Generate lots of OTP using multi processes and verity the uniqueness
func TestUniqueness(t *testing.T) {
	numberOfCodes := 50

	generatedCodes := map[string]struct{}{}

	for range numberOfCodes {
		code, err := Generate()

		assert.Nil(t, err, "Generate() unexpected error: %w", err)
		assert.NotContains(t, generatedCodes, code, "Generate() duplicated code detected: got %s twice", code)

		generatedCodes[code] = struct{}{}
	}

	assert.Len(t, generatedCodes, numberOfCodes)
}

func BenchmarkGenerate(b *testing.B) {
	for b.Loop() {
		if _, err := Generate(); err != nil {
			b.Fatal(err)
		}
	}
}
