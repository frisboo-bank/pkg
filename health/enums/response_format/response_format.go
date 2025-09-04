package responseformat

type (
	responseFormat int8
)

const (
	unknown responseFormat = iota // invalid unknown
	text
	json
)
