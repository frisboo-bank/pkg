package logrus

import (
	"os"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/logger/config"

	"github.com/nolleh/caption_json_formatter"
	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	level    logger.LogLevel
	encoding string
	logger   *logrus.Logger
	config   *config.LogOptions
}

var logrusLevelMapping = map[logger.LogLevel]logrus.Level{
	logger.LogLevelDebug: logrus.DebugLevel,
	logger.LogLevelInfo:  logrus.InfoLevel,
	logger.LogLevelWarn:  logrus.WarnLevel,
	logger.LogLevelError: logrus.ErrorLevel,
	logger.LogLevelPanic: logrus.PanicLevel,
	logger.LogLevelFatal: logrus.FatalLevel,
}

func NewLogrusLogger(cfg *config.LogOptions) logger.Logger {
	logger := &logrusLogger{
		level:  cfg.Level,
		config: cfg,
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
		logger.SetReportCaller(true)
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

func (l *logrusLogger) Debugw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) Error(v ...any) {
	l.logger.Error(v...)
}

func (l *logrusLogger) Errorf(format string, v ...any) {
	l.logger.Errorf(format, v...)
}

func (l *logrusLogger) Errorw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) Fatal(v ...any) {
	l.logger.Fatal(v...)
}

func (l *logrusLogger) Fatalf(format string, v ...any) {
	l.logger.Fatalf(format, v...)
}

func (l *logrusLogger) Fatalw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) Info(v ...any) {
	l.logger.Info(v...)
}

func (l *logrusLogger) Infof(format string, v ...any) {
	l.logger.Infof(format, v...)
}

func (l *logrusLogger) Infow(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) LogType() logger.LogType {
	return logger.LogTypeLogrus
}

func (l *logrusLogger) Panic(v ...any) {
	l.logger.Panic(v...)
}

func (l *logrusLogger) Panicf(format string, v ...any) {
	l.logger.Panicf(format, v...)
}

func (l *logrusLogger) Panicw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) Print(v ...any) {
	l.logger.Print(v...)
}

func (l *logrusLogger) Printf(format string, v ...any) {
	l.logger.Printf(format, v...)
}

func (l *logrusLogger) Printw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) Warn(v ...any) {
	l.logger.Warn(v...)
}

func (l *logrusLogger) Warnf(format string, v ...any) {
	l.logger.Warnf(format, v...)
}

func (l *logrusLogger) Warnw(message string, fields logger.Fields) {
	panic("unimplemented")
}

func (l *logrusLogger) WithName(name string) {
	l.logger.WithField(constants.LOGGER_NAME, name)
}
