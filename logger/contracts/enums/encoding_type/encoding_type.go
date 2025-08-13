package encodingtype

//go:generate goenums -f encoding_type.go

type (
	encodingType int8
)

const (
	unknown encodingType = iota
	json
	text
)
