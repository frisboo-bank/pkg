package noop

import (
	"frisboo-bank/pkg/logger/contracts"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

type NoopLogger struct {
	contracts.BaseLogger
}

func New() contracts.Logger {
	logger := &NoopLogger{}
	logger.BaseLogger.Init(logger)
	logger.SetupInstance()

	return logger
}

func (l *NoopLogger) SetupInstance() {}

func (l *NoopLogger) Log(level loglevel.LogLevel, v ...any) {}

func (l *NoopLogger) Logf(level loglevel.LogLevel, format string, v ...any) {}

func (l *NoopLogger) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {}

func (l *NoopLogger) Instance() any {
	return l
}

func (l *NoopLogger) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.NOOP
}
