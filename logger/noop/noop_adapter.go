package noop

import (
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

var _ contracts.LoggerAdapter = (*noopAdapter)(nil)

type noopAdapter struct{}

func New() contracts.LoggerAdapter {
	return &noopAdapter{}
}

func (n *noopAdapter) Setup(cfg *config.Config) error {
	return nil
}

func (n *noopAdapter) Log(level loglevel.LogLevel, v ...any) {}

func (n *noopAdapter) Logf(level loglevel.LogLevel, format string, v ...any) {}

func (n *noopAdapter) Logw(level loglevel.LogLevel, message string, fields contracts.Fields) {}

func (n *noopAdapter) Type() loggertype.LoggerType {
	return loggertype.LoggerTypes.NOOP
}
