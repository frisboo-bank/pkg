package logger

import (
	"frisboo-bank/pkg/logger/contracts"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
	"frisboo-bank/pkg/syserrors"
)

var _ contracts.Logger = (*logger)(nil)

type logger struct {
	adapter contracts.LoggerAdapter
}

func New(adapter contracts.LoggerAdapter) contracts.Logger {
	syserrors.Assert(adapter != nil, "adapter can't be nil")
	return &logger{adapter}
}

func (l *logger) Debug(v ...any)            { l.log(loglevel.LogLevels.DEBUGLEVEL, v...) }
func (l *logger) Debugf(f string, v ...any) { l.logf(loglevel.LogLevels.DEBUGLEVEL, f, v...) }
func (l *logger) Debugw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.DEBUGLEVEL, msg, fields)
}

func (l *logger) Error(v ...any)            { l.log(loglevel.LogLevels.ERRORLEVEL, v...) }
func (l *logger) Errorf(f string, v ...any) { l.logf(loglevel.LogLevels.ERRORLEVEL, f, v...) }
func (l *logger) Errorw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.ERRORLEVEL, msg, fields)
}

func (l *logger) Fatal(v ...any)            { l.log(loglevel.LogLevels.FATALLEVEL, v...) }
func (l *logger) Fatalf(f string, v ...any) { l.logf(loglevel.LogLevels.FATALLEVEL, f, v...) }
func (l *logger) Fatalw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.FATALLEVEL, msg, fields)
}

func (l *logger) Info(v ...any)            { l.log(loglevel.LogLevels.INFOLEVEL, v...) }
func (l *logger) Infof(f string, v ...any) { l.logf(loglevel.LogLevels.INFOLEVEL, f, v...) }
func (l *logger) Infow(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.INFOLEVEL, msg, fields)
}

func (l *logger) Panic(v ...any)            { l.log(loglevel.LogLevels.PANICLEVEL, v...) }
func (l *logger) Panicf(f string, v ...any) { l.logf(loglevel.LogLevels.PANICLEVEL, f, v...) }
func (l *logger) Panicw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.PANICLEVEL, msg, fields)
}

func (l *logger) Print(v ...any)            { l.log(loglevel.LogLevels.INFOLEVEL, v...) }
func (l *logger) Printf(f string, v ...any) { l.logf(loglevel.LogLevels.INFOLEVEL, f, v...) }
func (l *logger) Printw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.INFOLEVEL, msg, fields)
}

func (l *logger) Warn(v ...any)            { l.log(loglevel.LogLevels.WARNLEVEL, v...) }
func (l *logger) Warnf(f string, v ...any) { l.logf(loglevel.LogLevels.WARNLEVEL, f, v...) }
func (l *logger) Warnw(msg string, fields contracts.Fields) {
	l.logw(loglevel.LogLevels.WARNLEVEL, msg, fields)
}

func (l *logger) log(level loglevel.LogLevel, v ...any) { l.adapter.Log(level, v...) }
func (l *logger) logf(level loglevel.LogLevel, format string, v ...any) {
	l.adapter.Logf(level, format, v...)
}

func (l *logger) logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
	l.adapter.Logw(level, message, fields)
}

// Type implements contracts.Logger.
func (l *logger) Type() loggertype.LoggerType {
	return l.adapter.Type()
}
