package loglevel

type (
	logLevel int8
)

const (
	unknown logLevel = iota // unknown
	invalid                 // invalid
	debugLevel
	infoLevel
	warnLevel
	errorLevel
	panicLevel
	fatalLevel
	traceLevel
)
