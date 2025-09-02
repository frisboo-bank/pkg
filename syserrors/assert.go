package syserrors

import (
	"fmt"
	"os"
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
		nerr = New(err)
	default:
		fmt.Printf("assert err can only be an error or a string: get %q\n", err)
		os.Exit(1)
	}

	fmt.Println(Message(nerr, prefix))
	os.Exit(1)
}

func AssertNotNil(name string, value any) {
	Assert(value != nil, CantBeNilError(name))
}
