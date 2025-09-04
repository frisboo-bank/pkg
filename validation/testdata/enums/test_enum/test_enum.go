package testenum

type (
	testEnum int8
)

const (
	unknown testEnum = iota // unknown
	failed                  // invalid
	passed
	skipped
	scheduled
	running
	booked
)
