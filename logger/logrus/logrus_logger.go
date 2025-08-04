package logrus

import (
	"fmt"

	"frisboo-bank/pkg/logger/contracts"

	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

	"github.com/sirupsen/logrus"
)

type logrusLogger struct {
	contracts.BaseLogger
	instance *logrus.Logger
	// fields   contracts.Fields
}

func New() contracts.Logger {
	logger := &logrusLogger{
		instance: logrus.New(),
	}
	logger.BaseLogger.Init(logger)
	logger.SetupInstance()

	return logger
}

func (l *logrusLogger) SetupInstance() {
	if l.Output() != nil && l.instance.Out != l.Output() {
		l.instance.SetOutput(l.Output())
	}
}

func (l *logrusLogger) Log(level loglevel.LogLevel, v ...any) {
	if l.Prefix() != "" {
		v = append([]any{fmt.Sprintf("%s: ", l.Prefix())}, v...)
	}

	switch level {
	case loglevel.LogLevels.DEBUGLEVEL:
		l.instance.Debug(v...)
	case loglevel.LogLevels.ERRORLEVEL:
		l.instance.Error(v...)
	case loglevel.LogLevels.FATALLEVEL:
		l.instance.Fatal(v...)
	case loglevel.LogLevels.INFOLEVEL:
		l.instance.Info(v...)
	case loglevel.LogLevels.PANICLEVEL:
		l.instance.Panic(v...)
	case loglevel.LogLevels.WARNLEVEL:
		l.instance.Warn(v...)
	}
}

func (l *logrusLogger) Logf(level loglevel.LogLevel, format string, v ...any) {
	if l.Prefix() != "" {
		format = fmt.Sprintf("%s: %s", l.Prefix(), format)
	}

	switch level {
	case loglevel.LogLevels.DEBUGLEVEL:
		l.instance.Debugf(format, v...)
	case loglevel.LogLevels.ERRORLEVEL:
		l.instance.Errorf(format, v...)
	case loglevel.LogLevels.FATALLEVEL:
		l.instance.Fatalf(format, v...)
	case loglevel.LogLevels.INFOLEVEL:
		l.instance.Infof(format, v...)
	case loglevel.LogLevels.PANICLEVEL:
		l.instance.Panicf(format, v...)
	case loglevel.LogLevels.WARNLEVEL:
		l.instance.Warnf(format, v...)
	}
}

func (l *logrusLogger) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
	if l.Prefix() != "" {
		message = fmt.Sprintf("%s: %s", l.Prefix(), message)
	}
	entry := l.getEntry()

	switch level {
	case loglevel.LogLevels.DEBUGLEVEL:
		entry.Debug(message)
	case loglevel.LogLevels.ERRORLEVEL:
		entry.Error(message)
	case loglevel.LogLevels.FATALLEVEL:
		entry.Fatal(message)
	case loglevel.LogLevels.INFOLEVEL:
		entry.Info(message)
	case loglevel.LogLevels.PANICLEVEL:
		entry.Panic(message)
	case loglevel.LogLevels.WARNLEVEL:
		entry.Warn(message)
	}
}

func (l *logrusLogger) Instance() any {
	return l.instance
}

func (l *logrusLogger) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.LOGRUS
}

// // Helper to get entry with all contextual fields
func (l *logrusLogger) getEntry() *logrus.Entry {
	fields := logrus.Fields{}
	if l.Prefix() != "" {
		fields["prefix"] = l.Prefix
	}

	// if l.Name != "" {
	// 	fields[constants.LOGGER_NAME] = l.Name
	// }
	// maps.Copy(fields, l.Fields)

	return l.instance.WithFields(fields)
}
