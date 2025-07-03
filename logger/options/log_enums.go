package options

import "frisboo-bank/pkg/types/enum"

type (
	LogType      int32
	LogLevel     int32
	EncodingType int32
)

const (
	LogTypeLogrus LogType = iota
	LogTypeNoop
)

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelPanic
	LogLevelFatal
)

const (
	EncodingTypeJSON EncodingType = iota
	EncodingTypeText
)

var LogTypeEnum = enum.New[LogType](map[LogType]string{
	LogTypeLogrus: "logrus",
	LogTypeNoop:   "noop",
})

var LogLevelEnum = enum.New[LogLevel](map[LogLevel]string{
	LogLevelDebug: "debug",
	LogLevelInfo:  "info",
	LogLevelWarn:  "warn",
	LogLevelError: "error",
	LogLevelPanic: "panic",
	LogLevelFatal: "fatal",
})

var EncodingTypeEnum = enum.New[EncodingType](map[EncodingType]string{
	EncodingTypeJSON: "json",
	EncodingTypeText: "text",
})

func ParseLogType(name string) (LogType, error)           { return LogTypeEnum.FromName(name) }
func ParseLogLevel(name string) (LogLevel, error)         { return LogLevelEnum.FromName(name) }
func ParseEncodingType(name string) (EncodingType, error) { return EncodingTypeEnum.FromName(name) }

func LogTypeNames() []string      { return LogTypeEnum.Names() }
func LogLevelNames() []string     { return LogLevelEnum.Names() }
func EncodingTypeNames() []string { return EncodingTypeEnum.Names() }
