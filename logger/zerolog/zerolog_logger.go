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
	var evt *zerolog.Event
	switch level {
	case loglevel.LogLevels.DEBUG:
		evt = l.instance.Debug()
	case loglevel.LogLevels.INFO:
		evt = l.instance.Info()
	case loglevel.LogLevels.WARN:
		evt = l.instance.Warn()
	case loglevel.LogLevels.ERROR:
		evt = l.instance.Error()
	case loglevel.LogLevels.FATAL:
		evt = l.instance.Fatal()
	case loglevel.LogLevels.TRACE:
		evt = l.instance.Trace()
	default:
		evt = l.instance.Info()
	}
	for _, val := range v {
		evt.Msgf("%v", val)
	}
}

func (l *zerologLogger) Logf(level loglevel.LogLevel, format string, v ...any) {
	msg := format
	if len(v) > 0 {
		msg = fmt.Sprintf(format, v...)
	}
	switch level {
	case loglevel.LogLevels.DEBUG:
		l.instance.Debug().Msg(msg)
	case loglevel.LogLevels.INFO:
		l.instance.Info().Msg(msg)
	case loglevel.LogLevels.WARN:
		l.instance.Warn().Msg(msg)
	case loglevel.LogLevels.ERROR:
		l.instance.Error().Msg(msg)
	case loglevel.LogLevels.FATAL:
		l.instance.Fatal().Msg(msg)
	case loglevel.LogLevels.PANIC:
		l.instance.Panic().Msg(msg)
	default:
		l.instance.Info().Msg(msg)
	}
}

func (l *zerologLogger) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
}

func (l *zerologLogger) Instance() any {
	return l.instance
}

func (l *zerologLogger) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.ZEROLOG
}
