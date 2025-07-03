package noop

import (
	"frisboo-bank/pkg/logger/contracts"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

type noopLogger struct {
	prefix string
}

var _ contracts.Logger = (*noopLogger)(nil)

func NewNoopLogger() contracts.Logger {
	return newNoopLogger()
}

func newNoopLogger() contracts.Logger {
	return &noopLogger{}
}

func (n *noopLogger) Configure(cfg func(internalLoggerConfig any)) {
}

func (n *noopLogger) Debug(v ...any) {
}

func (n *noopLogger) Debugf(format string, v ...any) {
}

func (n *noopLogger) Debugw(message string, fields contracts.Fields) {
}

func (n *noopLogger) Error(v ...any) {
}

func (n *noopLogger) Errorf(format string, v ...any) {
}

func (n *noopLogger) Errorw(message string, fields contracts.Fields) {
}

func (n *noopLogger) Fatal(v ...any) {
}

func (n *noopLogger) Fatalf(format string, v ...any) {
}

func (n *noopLogger) Fatalw(message string, fields contracts.Fields) {
}

func (n *noopLogger) Info(v ...any) {
}

func (n *noopLogger) Infof(format string, v ...any) {
}

func (n *noopLogger) Infow(message string, fields contracts.Fields) {
}

func (n *noopLogger) LogType() logtype.LogType {
	return logtype.LogTypes.NOOP
}

func (n *noopLogger) Panic(v ...any) {
}

func (n *noopLogger) Panicf(format string, v ...any) {
}

func (n *noopLogger) Panicw(message string, fields contracts.Fields) {
}

func (n *noopLogger) Print(v ...any) {
}

func (n *noopLogger) Printf(format string, v ...any) {
}

func (n *noopLogger) Printw(message string, fields contracts.Fields) {
}

func (n *noopLogger) Warn(v ...any) {
}

func (n *noopLogger) Warnf(format string, v ...any) {
}

func (n *noopLogger) Warnw(message string, fields contracts.Fields) {
}

func (n *noopLogger) WithName(name string) contracts.Logger {
	return n
}

// WithPrefix set the prefix
func (n *noopLogger) WithPrefix(prefix string) contracts.Logger {
	return n
}

// GetPrefix get the prefix
func (n *noopLogger) GetPrefix() string {
	return n.prefix
}
