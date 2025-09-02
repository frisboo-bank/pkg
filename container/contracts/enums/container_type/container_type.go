package containertype

//go:generate goenums -f -c container_type.go

type (
	containerType int8
)

const (
	unknown containerType = iota // invalid unknown
	dig
)
