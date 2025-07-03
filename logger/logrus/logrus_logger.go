package logrus

import (
	"maps"
	"os"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/options"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
	"frisboo-bank/pkg/logger/utils"

	"github.com/sirupsen/logrus"
)

var logrusLevelMapping = map[loglevel.LogLevel]logrus.Level{
	loglevel.LogLevels.DEBUG_LEVEL: logrus.DebugLevel,
	loglevel.LogLevels.INFO_LEVEL:  logrus.InfoLevel,
	loglevel.LogLevels.WARN_LEVEL:  logrus.WarnLevel,
	loglevel.LogLevels.ERROR_LEVEL: logrus.ErrorLevel,
	loglevel.LogLevels.PANIC_LEVEL: logrus.PanicLevel,
	loglevel.LogLevels.FATAL_LEVEL: logrus.FatalLevel,
}

type logrusLogger struct {
	level  loglevel.LogLevel
	logger *logrus.Logger
	config *options.LogOptions
	prefix string
	name   string
	fields contracts.Fields
}

var _ contracts.Logger = (*logrusLogger)(nil)

func NewLogrusLogger(config *options.LogOptions) contracts.Logger {
	return newLogrusLogger(config)
}

func newLogrusLogger(config *options.LogOptions) contracts.Logger {
	logger := &logrusLogger{
		level:  config.Level,
		config: config,
		fields: make(contracts.Fields),
	}
	logger.initLogger()
	return logger
}

func (l *logrusLogger) initLogger() {
	logger := logrus.New()
	logger.SetLevel(l.GetLogLevel())
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(false)

	switch l.config.Encoding {
	default: // Default to text
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		})
	}

	l.logger = logger
}

func (l *logrusLogger) GetLogLevel() logrus.Level {
	level, exist := logrusLevelMapping[l.level]
	if !exist {
		return logrus.ErrorLevel
	}
	return level
}

func (l *logrusLogger) Configure(cfg func(internalLoggerConfig any)) {
	cfg(l.logger)
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
	for k, v := range l.fields {
		fields[k] = v
	}

	return l.logger.WithFields(fields)
}

// --- Standard logging methods ---
func (l *logrusLogger) log(level logrus.Level, v ...any) {
	msg := utils.ConcatPrefix(l.prefix, v...)
	switch level {
	case logrus.DebugLevel:
		l.logger.Debug(msg...)
	case logrus.InfoLevel:
		l.logger.Info(msg...)
	case logrus.WarnLevel:
		l.logger.Warn(msg...)
	case logrus.ErrorLevel:
		l.logger.Error(msg...)
	case logrus.FatalLevel:
		l.logger.Fatal(msg...)
	case logrus.PanicLevel:
		l.logger.Panic(msg...)
	}
}

func (l *logrusLogger) logf(level logrus.Level, format string, v ...any) {
	format, v = utils.ConcatPrefixf(l.prefix, format, v...)
	switch level {
	case logrus.DebugLevel:
		l.logger.Debugf(format, v...)
	case logrus.InfoLevel:
		l.logger.Infof(format, v...)
	case logrus.WarnLevel:
		l.logger.Warnf(format, v...)
	case logrus.ErrorLevel:
		l.logger.Errorf(format, v...)
	case logrus.FatalLevel:
		l.logger.Fatalf(format, v...)
	case logrus.PanicLevel:
		l.logger.Panicf(format, v...)
	}
}

func (l *logrusLogger) logw(level logrus.Level, message string, fields contracts.Fields) {
	message = utils.ConcatPrefixStr(l.prefix, message)
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

// --- Concrete implementations ---
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

func (l *logrusLogger) WithName(name string) contracts.Logger {
	newLogger := *l
	newLogger.name = name
	return &newLogger
}

func (l *logrusLogger) WithPrefix(prefix string) contracts.Logger {
	newLogger := *l
	newLogger.prefix = prefix
	return &newLogger
}

func (l *logrusLogger) WithFields(fields contracts.Fields) contracts.Logger {
	newLogger := *l
	newLogger.fields = make(contracts.Fields, len(l.fields)+len(fields))
	maps.Copy(newLogger.fields, l.fields)
	maps.Copy(newLogger.fields, fields)

	return &newLogger
}

func (l *logrusLogger) GetPrefix() string {
	return l.prefix
}
