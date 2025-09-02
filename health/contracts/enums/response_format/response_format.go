package responseformat

//go:generate goenums -f -c response_format.go

type (
	responseFormat int8
)

const (
	unknown responseFormat = iota // invalid unknown
	text
	json
)
