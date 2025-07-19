package contracts

import (
	"io"

	"frisboo-bank/pkg/logger/options"
	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

type (
	Fields map[string]any

	Logger interface {
		WithOptions(options *options.LogOptions) Logger
		WithCaller(withCaller bool, depth int) Logger
		WithEncoding(encoding encodingtype.EncodingType) Logger
		WithLevel(logLevel loglevel.LogLevel) Logger
		WithName(name string) Logger
		WithOutput(output io.Writer) Logger
		WithPrefix(prefix string) Logger
		WithTracer(witlTracer bool) Logger

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
		Panic(v ...any)
		Panicf(format string, v ...any)
		Panicw(message string, fields Fields)
		Print(v ...any)
		Printf(format string, v ...any)
		Printw(message string, fields Fields)
		Warn(v ...any)
		Warnf(format string, v ...any)
		Warnw(message string, fields Fields)

		LogType() logtype.LogType
		Instance() any
	}
)
