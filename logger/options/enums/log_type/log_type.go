package logtype

//go:generate goenums -f log_type.go

type (
	logType int8
)

const (
	unknown logType = iota // invalid unknown
	logrus                 // logrus
	noop                   // noop
)
