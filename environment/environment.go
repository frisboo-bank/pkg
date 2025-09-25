package environment

type (
	environment int8
)

const (
	unknown environment = iota // invalid
	development
	preprod
	production
	testing
)
