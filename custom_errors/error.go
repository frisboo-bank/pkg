package customErrors

import "errors"

type fatalError struct{ error }

func (e *fatalError) Error() string { return e.error.Error() }
func (e *fatalError) Unwrap() error { return e.error }

func WrapAsFatal(err error) error {
	return &fatalError{err}
}

func IsFatal(err error) bool {
	var expectedError *fatalError
	return errors.As(err, &expectedError)
}

func WrapWith(sentinel error, err error) error {
	return errors.Join(sentinel, err)
}
