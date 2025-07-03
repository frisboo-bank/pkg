package httpservertype

//go:generate goenums -f -c http_server_type.go

type (
	httpServerType int8
)

const (
	unknown httpServerType = iota // invalid unknown
	gin                           // gin
)
