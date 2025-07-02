package logrus

import (
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/options"
	"os"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
)

var logrusLevelMapping = map[options.LogLevel]logrus.Level{
	options.LogLevelDebug: logrus.DebugLevel,
	options.LogLevelInfo:  logrus.InfoLevel,
	options.LogLevelWarn:  logrus.WarnLevel,
	options.LogLevelError: logrus.ErrorLevel,
	options.LogLevelPanic: logrus.PanicLevel,
	options.LogLevelFatal: logrus.FatalLevel,
}

type logrusLogger struct {
	level    options.LogLevel
	encoding string
	logger   *logrus.Logger
	config   *options.LogOptions
}

var _ contracts.Logger = (*logrusLogger)(nil)

func NewLogrusLogger(config *options.LogOptions) contracts.Logger {
	return newLogrusLogger(config)
}

func newLogrusLogger(config *options.LogOptions) contracts.Logger {
	logger := &logrusLogger{
		level:  config.Level,
		config: config,
	}
	logger.initLogger()

	return logger
}

func (l *logrusLogger) initLogger() {
	logLevel := l.GetLogLevel()

	logger := logrus.New()
	logger.SetLevel(logLevel)
	logger.SetOutput(os.Stdout)

	switch true {
	case true:
		logger.SetReportCaller(false)
		logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
		})
	default:
		logger.SetReportCaller(false)
		logger.SetFormatter(&caption_json_formatter.Formatter{PrettyPrint: true})
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

func (l *logrusLogger) Debug(v ...any) {
	l.logger.Debug(v...)
}

func (l *logrusLogger) Debugf(format string, v ...any) {
	l.logger.Debugf(format, v...)
}

func (l *logrusLogger) Debugw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Debug(message)
}

func (l *logrusLogger) Error(v ...any) {
	l.logger.Error(v...)
}

func (l *logrusLogger) Errorf(format string, v ...any) {
	l.logger.Errorf(format, v...)
}

func (l *logrusLogger) Errorw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Error(message)
}

func (l *logrusLogger) Fatal(v ...any) {
	l.logger.Fatal(v...)
}

func (l *logrusLogger) Fatalf(format string, v ...any) {
	l.logger.Fatalf(format, v...)
}

func (l *logrusLogger) Fatalw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Fatal(message)
}

func (l *logrusLogger) Info(v ...any) {
	l.logger.Info(v...)
}

func (l *logrusLogger) Infof(format string, v ...any) {
	l.logger.Infof(format, v...)
}

func (l *logrusLogger) Infow(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Info(message)
}

func (l *logrusLogger) LogType() options.LogType {
	return options.TypeLogrus
}

func (l *logrusLogger) Panic(v ...any) {
	l.logger.Panic(v...)
}

func (l *logrusLogger) Panicf(format string, v ...any) {
	l.logger.Panicf(format, v...)
}

func (l *logrusLogger) Panicw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Panic(message)
}

func (l *logrusLogger) Print(v ...any) {
	l.logger.Print(v...)
}

func (l *logrusLogger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *logrusLogger) Printw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Print(message)
}

func (l *logrusLogger) Warn(v ...any) {
	l.logger.Warn(v...)
}

func (l *logrusLogger) Warnf(format string, v ...any) {
	l.logger.Warnf(format, v...)
}

func (l *logrusLogger) Warnw(message string, fields contracts.Fields) {
	entry := l.mapToFields(fields)
	entry.Warn(message)
}

func (l *logrusLogger) WithName(name string) contracts.Logger {
	l.logger.WithField(constants.LOGGER_NAME, name)
	return l
}

func (l *logrusLogger) mapToFields(fields map[string]any) *logrus.Entry {
	return l.logger.WithFields(logrus.Fields{"ss": 1})
}
