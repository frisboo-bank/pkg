package otp

import (
	"crypto/rand"
	"math/big"

	"frisboo-bank/pkg/syserrors"
)

const (
	Charset       = "0123456789"
	DefaultLength = 6
)

func Generate() (string, error) {
	return GenerateWithLength(DefaultLength)
}

func GenerateWithLength(length int) (string, error) {
	if length <= 0 {
		return "", syserrors.Newf("length out of range: got %d", length)
	}

	code := make([]byte, length)

	for i := range length {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(Charset))))
		if err != nil {
			return "", syserrors.Wrap(err, "failed to generate random index")
		}

		code[i] = Charset[index.Int64()]
	}

	return string(code), nil
}
