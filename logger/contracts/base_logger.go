package contracts

import (
	"io"

	"frisboo-bank/pkg/logger/config"

	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
)

type BaseLogger struct {
	callDepth     int
	callerEnabled bool
	encoding      encodingtype.EncodingType
	level         loglevel.LogLevel
	name          string
	output        io.Writer
	prefix        string
	tracerEnabled bool

	internal LoggerInternal
}

var _ loggerConfig = (*BaseLogger)(nil)

func (b *BaseLogger) Init(internal LoggerInternal) {
	b.internal = internal
}

func (b *BaseLogger) WithConfig(cfg *config.LoggerConfig) Logger {
	b.callDepth = cfg.CallDepth
	b.callerEnabled = cfg.CallerEnabled
	b.encoding = cfg.Encoding
	b.level = cfg.Level
	b.tracerEnabled = cfg.EnableTracing

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithCallDepth(callDepth int) Logger {
	b.callDepth = callDepth

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithCaller(withCaller bool) Logger {
	b.callerEnabled = withCaller

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithEncoding(encoding encodingtype.EncodingType) Logger {
	b.encoding = encoding

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithLevel(logLevel loglevel.LogLevel) Logger {
	b.level = logLevel

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithName(name string) Logger {
	b.name = name

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithOutput(output io.Writer) Logger {
	b.output = output

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithPrefix(prefix string) Logger {
	b.prefix = prefix

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) WithTracer(withTracer bool) Logger {
	b.tracerEnabled = withTracer

	b.internal.SetupInstance()
	return b.internal.(Logger)
}

func (b *BaseLogger) CallDepth() int                      { return b.callDepth }
func (b *BaseLogger) CallerEnabled() bool                 { return b.callerEnabled }
func (b *BaseLogger) Encoding() encodingtype.EncodingType { return b.encoding }
func (b *BaseLogger) Level() loglevel.LogLevel            { return b.level }
func (b *BaseLogger) Name() string                        { return b.name }
func (b *BaseLogger) Output() io.Writer                   { return b.output }
func (b *BaseLogger) Prefix() string                      { return b.prefix }
func (b *BaseLogger) TracerEnabled() bool                 { return b.tracerEnabled }

func (b *BaseLogger) Debug(v ...any) { b.internal.Log(loglevel.LogLevels.DEBUGLEVEL, v...) }

func (b *BaseLogger) Debugf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.DEBUGLEVEL, format, v...)
}

func (b *BaseLogger) Debugw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.DEBUGLEVEL, message, fields)
}

func (b *BaseLogger) Error(v ...any) { b.internal.Log(loglevel.LogLevels.ERRORLEVEL, v...) }

func (b *BaseLogger) Errorf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.ERRORLEVEL, format, v...)
}

func (b *BaseLogger) Errorw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.ERRORLEVEL, message, fields)
}

func (b *BaseLogger) Fatal(v ...any) { b.internal.Log(loglevel.LogLevels.FATALLEVEL, v...) }

func (b *BaseLogger) Fatalf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.FATALLEVEL, format, v...)
}

func (b *BaseLogger) Fatalw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.FATALLEVEL, message, fields)
}

func (b *BaseLogger) Info(v ...any) { b.internal.Log(loglevel.LogLevels.INFOLEVEL, v...) }

func (b *BaseLogger) Infof(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.INFOLEVEL, format, v...)
}

func (b *BaseLogger) Infow(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.INFOLEVEL, message, fields)
}

func (b *BaseLogger) Panic(v ...any) { b.internal.Log(loglevel.LogLevels.PANICLEVEL, v...) }

func (b *BaseLogger) Panicf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.PANICLEVEL, format, v...)
}

func (b *BaseLogger) Panicw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.PANICLEVEL, message, fields)
}

func (b *BaseLogger) Print(v ...any) { b.internal.Log(loglevel.LogLevels.INFOLEVEL, v...) }

func (b *BaseLogger) Printf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.INFOLEVEL, format, v...)
}

func (b *BaseLogger) Printw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.INFOLEVEL, message, fields)
}

func (b *BaseLogger) Warn(v ...any) { b.internal.Log(loglevel.LogLevels.WARNLEVEL, v...) }

func (b *BaseLogger) Warnf(format string, v ...any) {
	b.internal.Logf(loglevel.LogLevels.WARNLEVEL, format, v...)
}

func (b *BaseLogger) Warnw(message string, fields Fields) {
	b.internal.Logw(loglevel.LogLevels.WARNLEVEL, message, fields)
}
