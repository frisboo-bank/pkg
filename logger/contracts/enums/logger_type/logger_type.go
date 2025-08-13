package loggertype

//go:generate goenums -f logger_type.go

type (
	loggerType int8
)

const (
	unknown loggerType = iota
	logrus
	noop
	zerolog
)
