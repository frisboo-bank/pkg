package encodingtype

//go:generate goenums -f encoding_type.go

type (
	encodingType int8
)

const (
	unknown encodingType = iota // invalid unknown
	json                        // json
	text                        // text
)
