package zerolog

import (
	"fmt"

	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"

	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

	"github.com/rs/zerolog"
)

var _ contracts.LoggerAdapter = (*zerologAdapter)(nil)

var loglevelsMapping = map[loglevel.LogLevel]zerolog.Level{
	loglevel.LogLevels.DEBUGLEVEL: zerolog.DebugLevel,
	loglevel.LogLevels.INFOLEVEL:  zerolog.InfoLevel,
	loglevel.LogLevels.WARNLEVEL:  zerolog.WarnLevel,
	loglevel.LogLevels.ERRORLEVEL: zerolog.ErrorLevel,
	loglevel.LogLevels.FATALLEVEL: zerolog.FatalLevel,
	loglevel.LogLevels.PANICLEVEL: zerolog.PanicLevel,
}

type zerologAdapter struct {
	cfg    *config.Config
	logger zerolog.Logger
}

func New(cfg *config.Config) contracts.LoggerAdapter {
	logger := zerolog.New(cfg.Output)

	return &zerologAdapter{
		cfg:    cfg,
		logger: logger,
	}
}

func (z *zerologAdapter) Log(level loglevel.LogLevel, v ...any) {
	event := z.logger.WithLevel(mapToZeroLogLevel(level))
	event.Msg(fmt.Sprint(v...))
}

func (z *zerologAdapter) Logf(level loglevel.LogLevel, format string, v ...any) {
	event := z.logger.WithLevel(mapToZeroLogLevel(level))
	event.Msgf(format, v...)
}

func (z *zerologAdapter) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
	event := z.logger.
		WithLevel(mapToZeroLogLevel(level)).
		Fields(fields)
	event.Msg(message)
}

func (z *zerologAdapter) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.ZEROLOG
}

func mapToZeroLogLevel(level loglevel.LogLevel) zerolog.Level {
	if lv, ok := loglevelsMapping[level]; ok {
		return lv
	}
	return zerolog.InfoLevel
}
