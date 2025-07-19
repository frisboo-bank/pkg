package logrus

import (
	"fmt"
	"io"
	"maps"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/options"
	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
	"github.com/uptrace/opentelemetry-go-extra/otellogrus"
)

var logrusLevelMapping = map[loglevel.LogLevel]logrus.Level{
	loglevel.LogLevels.DEBUG_LEVEL: logrus.DebugLevel,
	loglevel.LogLevels.ERROR_LEVEL: logrus.ErrorLevel,
	loglevel.LogLevels.FATAL_LEVEL: logrus.FatalLevel,
	loglevel.LogLevels.INFO_LEVEL:  logrus.InfoLevel,
	loglevel.LogLevels.PANIC_LEVEL: logrus.PanicLevel,
	loglevel.LogLevels.WARN_LEVEL:  logrus.WarnLevel,
}

type logrusLogger struct {
	callDepth     int
	callerEnabled bool
	enableTracing bool
	encoding      encodingtype.EncodingType
	level         loglevel.LogLevel
	name          string
	output        io.Writer
	prefix        string

	instance *logrus.Logger
	fields   contracts.Fields
}

var _ contracts.Logger = (*logrusLogger)(nil)

func (l *logrusLogger) WithOptions(options *options.LogOptions) contracts.Logger {
	return l.
		WithCaller(options.CallerEnabled, options.CallDepth).
		WithEncoding(options.Encoding).
		WithLevel(options.Level).
		WithTracer(options.EnableTracing)
}

func (l *logrusLogger) WithCaller(enabled bool, depth int) contracts.Logger {
	l.callerEnabled = enabled
	l.callDepth = depth

	l.instance.SetReportCaller(enabled)

	return l
}

func (l *logrusLogger) WithEncoding(encoding encodingtype.EncodingType) contracts.Logger {
	l.encoding = encoding

	switch encoding {
	case encodingtype.EncodingTypes.JSON:
		l.instance.SetFormatter(&caption_json_formatter.Formatter{
			PrettyPrint: true,
		})
	default:
		l.instance.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		})
	}

	return l
}

func (l *logrusLogger) WithLevel(logLevel loglevel.LogLevel) contracts.Logger {
	l.level = logLevel

	l.instance.SetLevel(l.getLogLevel())

	return l
}

func (l *logrusLogger) WithName(name string) contracts.Logger {
	l.name = name
	return l
}

func (l *logrusLogger) WithOutput(output io.Writer) contracts.Logger {
	l.output = output

	l.instance.SetOutput(output)

	return l
}

func (l *logrusLogger) WithPrefix(prefix string) contracts.Logger {
	l.prefix = prefix
	return l
}

func (l *logrusLogger) WithTracer(withTracer bool) contracts.Logger {
	l.enableTracing = withTracer

	if l.enableTracing {
		l.instance.Debug("Tracing enabled. Current hooks: ", l.instance.Hooks)

		l.instance.AddHook(otellogrus.NewHook(otellogrus.WithLevels(
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		)))
	}

	return l
}

func NewLogrusLogger() contracts.Logger {
	logger := &logrusLogger{
		callerEnabled: false,
		enableTracing: false,
		encoding:      options.Encoding,
		fields:        make(contracts.Fields),
		level:         options.Level,
		output:        options.Output,
		instance:      logrus.New(),
	}

	logger.initLogger()

	return logger
}

func (l *logrusLogger) initLogger() {
	l.WithCaller(l.callerEnabled, l.callDepth).
		WithEncoding(l.encoding).
		WithLevel(l.level).
		WithName(l.name).
		WithOutput(l.output).
		WithPrefix(l.prefix).
		WithTracer(l.enableTracing)
}

func (l *logrusLogger) Debug(v ...any) {
	l.log(logrus.DebugLevel, v...)
}

func (l *logrusLogger) Debugf(format string, v ...any) {
	l.logf(logrus.DebugLevel, format, v...)
}

func (l *logrusLogger) Debugw(message string, fields contracts.Fields) {
	l.logw(logrus.DebugLevel, message, fields)
}

func (l *logrusLogger) Info(v ...any) {
	l.log(logrus.InfoLevel, v...)
}

func (l *logrusLogger) Infof(format string, v ...any) {
	l.logf(logrus.InfoLevel, format, v...)
}

func (l *logrusLogger) Infow(message string, fields contracts.Fields) {
	l.logw(logrus.InfoLevel, message, fields)
}

func (l *logrusLogger) Warn(v ...any) {
	l.log(logrus.WarnLevel, v...)
}

func (l *logrusLogger) Warnf(format string, v ...any) {
	l.logf(logrus.WarnLevel, format, v...)
}

func (l *logrusLogger) Warnw(message string, fields contracts.Fields) {
	l.logw(logrus.WarnLevel, message, fields)
}

func (l *logrusLogger) Error(v ...any) {
	l.log(logrus.ErrorLevel, v...)
}

func (l *logrusLogger) Errorf(format string, v ...any) {
	l.logf(logrus.ErrorLevel, format, v...)
}

func (l *logrusLogger) Errorw(message string, fields contracts.Fields) {
	l.logw(logrus.ErrorLevel, message, fields)
}

func (l *logrusLogger) Fatal(v ...any) {
	l.log(logrus.FatalLevel, v...)
}

func (l *logrusLogger) Fatalf(format string, v ...any) {
	l.logf(logrus.FatalLevel, format, v...)
}

func (l *logrusLogger) Fatalw(message string, fields contracts.Fields) {
	l.logw(logrus.FatalLevel, message, fields)
}

func (l *logrusLogger) Panic(v ...any) {
	l.log(logrus.PanicLevel, v...)
}

func (l *logrusLogger) Panicf(format string, v ...any) {
	l.logf(logrus.PanicLevel, format, v...)
}

func (l *logrusLogger) Panicw(message string, fields contracts.Fields) {
	l.logw(logrus.PanicLevel, message, fields)
}

func (l *logrusLogger) Print(v ...any) {
	l.log(logrus.InfoLevel, v...)
}

func (l *logrusLogger) Printf(format string, v ...any) {
	l.logf(logrus.InfoLevel, format, v...)
}

func (l *logrusLogger) Printw(message string, fields contracts.Fields) {
	l.logw(logrus.InfoLevel, message, fields)
}

func (l *logrusLogger) LogType() logtype.LogType {
	return logtype.LogTypes.LOGRUS
}

func (l *logrusLogger) Instance() any {
	return l.instance
}

// Helper to get the log level in logrus format
func (l *logrusLogger) getLogLevel() logrus.Level {
	level, exist := logrusLevelMapping[l.level]
	if !exist {
		return logrus.ErrorLevel
	}
	return level
}

// Helper to get entry with all contextual fields
func (l *logrusLogger) getEntry() *logrus.Entry {
	fields := logrus.Fields{}
	if l.prefix != "" {
		fields["prefix"] = l.prefix
	}

	if l.name != "" {
		fields[constants.LOGGER_NAME] = l.name
	}
	maps.Copy(fields, l.fields)

	return l.instance.WithFields(fields)
}

// --- Standard logging methods ---
func (l *logrusLogger) log(level logrus.Level, v ...any) {
	if l.prefix != "" {
		v = append([]any{fmt.Sprintf("%s: ", l.prefix)}, v...)
	}

	switch level {
	case logrus.DebugLevel:
		l.instance.Debug(v...)
	case logrus.InfoLevel:
		l.instance.Info(v...)
	case logrus.WarnLevel:
		l.instance.Warn(v...)
	case logrus.ErrorLevel:
		l.instance.Error(v...)
	case logrus.FatalLevel:
		l.instance.Fatal(v...)
	case logrus.PanicLevel:
		l.instance.Panic(v...)
	}
}

func (l *logrusLogger) logf(level logrus.Level, format string, v ...any) {
	if l.prefix != "" {
		format = fmt.Sprintf("%s: %s", l.prefix, format)
	}

	switch level {
	case logrus.DebugLevel:
		l.instance.Debugf(format, v...)
	case logrus.InfoLevel:
		l.instance.Infof(format, v...)
	case logrus.WarnLevel:
		l.instance.Warnf(format, v...)
	case logrus.ErrorLevel:
		l.instance.Errorf(format, v...)
	case logrus.FatalLevel:
		l.instance.Fatalf(format, v...)
	case logrus.PanicLevel:
		l.instance.Panicf(format, v...)
	}
}

func (l *logrusLogger) logw(level logrus.Level, message string, fields contracts.Fields) {
	if l.prefix != "" {
		message = fmt.Sprintf("%s: %s", l.prefix, message)
	}
	entry := l.getEntry()

	switch level {
	case logrus.DebugLevel:
		entry.Debug(message)
	case logrus.InfoLevel:
		entry.Info(message)
	case logrus.WarnLevel:
		entry.Warn(message)
	case logrus.ErrorLevel:
		entry.Error(message)
	case logrus.FatalLevel:
		entry.Fatal(message)
	case logrus.PanicLevel:
		entry.Panic(message)
	}
}
