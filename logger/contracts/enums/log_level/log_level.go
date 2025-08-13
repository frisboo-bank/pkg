package loglevel

//go:generate goenums -f log_level.go

type (
	logLevel int8
)

const (
	unknownLevel logLevel = iota
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	panicLevel
	fatalLevel
	traceLevel
)
