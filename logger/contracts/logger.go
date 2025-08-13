package contracts

import (
	"frisboo-bank/pkg/logger/config"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

type (
	Fields map[string]any

	Logger interface {
		Debug(v ...any)
		Debugf(format string, v ...any)
		Debugw(message string, fields Fields)
		Error(v ...any)
		Errorf(format string, v ...any)
		Errorw(message string, fields Fields)
		Fatal(v ...any)
		Fatalf(format string, v ...any)
		Fatalw(message string, fields Fields)
		Info(v ...any)
		Infof(format string, v ...any)
		Infow(message string, fields Fields)
		Panic(v ...any)
		Panicf(format string, v ...any)
		Panicw(message string, fields Fields)
		Print(v ...any)
		Printf(format string, v ...any)
		Printw(message string, fields Fields)
		Type() loggertype.LoggerType
		Warn(v ...any)
		Warnf(format string, v ...any)
		Warnw(message string, fields Fields)
	}

	LoggerAdapter interface {
		Log(level loglevel.LogLevel, v ...any)
		Logf(level loglevel.LogLevel, format string, v ...any)
		Logw(level loglevel.LogLevel, message string, fields Fields)
		Setup(cfg *config.Config) error
		Type() loggertype.LoggerType
	}
)
