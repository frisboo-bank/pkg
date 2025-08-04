package zerolog

import (
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

	"github.com/rs/zerolog"
)

type zerologLogger struct {
	contracts.BaseLogger
	instance zerolog.Logger
}

func New() contracts.Logger {
	logger := &zerologLogger{}
	logger.Init(logger)
	logger.SetupInstance()

	return logger
}

func (l *zerologLogger) SetupInstance() {
	output := l.Output()
	if output == nil {
		output = config.Output
	}

	l.instance = zerolog.New(output).With().Timestamp().Logger()
}

func (l *zerologLogger) Log(level loglevel.LogLevel, v ...any) {
}

func (l *zerologLogger) Logf(level loglevel.LogLevel, format string, v ...any) {
}

func (l *zerologLogger) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
}

func (l *zerologLogger) Instance() any {
	return l.instance
}

func (l *zerologLogger) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.ZEROLOG
}
