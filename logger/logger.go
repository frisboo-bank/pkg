package logger

import "errors"

type (
	Fields   map[string]any
	LogType  string
	LogLevel string
)

var ErrLogger = errors.New("Logger: something wrong happened")

const (
	LogTypeLogrus = LogType("logrus")
)

const (
	LogLevelDebug = LogLevel("debug")
	LogLevelInfo  = LogLevel("info")
	LogLevelWarn  = LogLevel("warn")
	LogLevelError = LogLevel("error")
	LogLevelPanic = LogLevel("panic")
	LogLevelFatal = LogLevel("fatal")
)

type Logger interface {
	Configure(cfg func(internalLoggerConfig any))
	Debug(v ...any)
	Debugf(format string, v ...any)
	Debugw(message string, fields Fields)
	Error(v ...any)
	Errorf(format string, v ...any)
	Errorw(message string, fields Fields)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Fatalw(message string, fields Fields)
	Info(v ...any)
	Infof(format string, v ...any)
	Infow(message string, fields Fields)
	LogType() LogType
	Panic(v ...any)
	Panicf(format string, v ...any)
	Panicw(message string, fields Fields)
	Print(v ...any)
	Printf(format string, v ...any)
	Printw(message string, fields Fields)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Warnw(message string, fields Fields)
	WithName(name string)
}
