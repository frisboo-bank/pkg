package validation

import (
	"fmt"
	"os"
	"strings"

	"frisboo-bank/pkg/syserrors"
)

func Assert(condition bool, err any, prefix ...string) {
	if condition {
		return
	}

	var nerr error
	switch err := err.(type) {
	case error:
		nerr = err
	case string:
		nerr = syserrors.New(err)
	default:
		fmt.Printf("assert err can only be an error or a string: get %v\n", err)
		os.Exit(1)
	}

	fmt.Println(syserrors.Message(nerr, prefix))
	os.Exit(1)
}

func AssertNotNil(name string, value any) {
	Assert(value != nil, syserrors.CantBeNilError(name))
}

func AssertNotEmpty(name string, value string) {
	Assert(strings.TrimSpace(value) != "", syserrors.CantBeEmptyError(name))
}
