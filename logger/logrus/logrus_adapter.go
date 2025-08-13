package logrus

import (
	"maps"
	"time"

	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"

	"github.com/sirupsen/logrus"
)

var _ contracts.LoggerAdapter = (*logrusAdapter)(nil)

var loglevelsMapping = map[loglevel.LogLevel]logrus.Level{
	loglevel.LogLevels.DEBUGLEVEL: logrus.DebugLevel,
	loglevel.LogLevels.INFOLEVEL:  logrus.InfoLevel,
	loglevel.LogLevels.WARNLEVEL:  logrus.WarnLevel,
	loglevel.LogLevels.ERRORLEVEL: logrus.ErrorLevel,
	loglevel.LogLevels.FATALLEVEL: logrus.FatalLevel,
	loglevel.LogLevels.PANICLEVEL: logrus.PanicLevel,
	loglevel.LogLevels.TRACELEVEL: logrus.TraceLevel,
}

type logrusAdapter struct {
	cfg    *config.Config
	logger *logrus.Logger
}

func New() contracts.LoggerAdapter {
	return &logrusAdapter{}
}

func (l *logrusAdapter) Setup(cfg *config.Config) error {
	logger := logrus.New()
	logger.SetLevel(mapToLogrusLevel(cfg.Level))
	logger.SetOutput(cfg.Output)
	logger.SetReportCaller(cfg.CallerEnabled)

	if cfg.Prefix != "" {
		logger.AddHook(&prefixHook{
			prefix: cfg.Prefix,
		})
	}

	switch cfg.Encoding {
	case encodingtype.EncodingTypes.JSON:
		logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	case encodingtype.EncodingTypes.TEXT:
		fallthrough
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339Nano,
		})
	}

	l.cfg = cfg
	l.logger = logger

	return nil
}

func (l *logrusAdapter) Log(level loglevel.LogLevel, v ...any) {
	l.logger.Log(mapToLogrusLevel(level), v...)
}

func (l *logrusAdapter) Logf(level loglevel.LogLevel, format string, v ...any) {
	l.logger.Logf(mapToLogrusLevel(level), format, v...)
}

func (l *logrusAdapter) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {
	lFields := make(contracts.Fields, len(fields))
	maps.Copy(lFields, fields)

	l.logger.
		WithFields(logrus.Fields(lFields)).
		Log(mapToLogrusLevel(level), message)
}

func (l *logrusAdapter) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.LOGRUS
}

func mapToLogrusLevel(level loglevel.LogLevel) logrus.Level {
	if lv, ok := loglevelsMapping[level]; ok {
		return lv
	}
	return logrus.InfoLevel
}
