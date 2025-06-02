package utils

func Assert(condition bool, err any) {
	if !condition {
		panic(err)
	}
}
