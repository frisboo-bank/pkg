package noop

import (
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/options"
	"io"

	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

type noopLogger struct{}

var _ contracts.Logger = (*noopLogger)(nil)

func (n *noopLogger) WithOptions(options *options.LogOptions) contracts.Logger {
	return n
}

func (n *noopLogger) WithCaller(_ bool, _ int) contracts.Logger {
	return n
}

func (n *noopLogger) WithEncoding(_ encodingtype.EncodingType) contracts.Logger { return n }

func (n *noopLogger) WithLevel(_ loglevel.LogLevel) contracts.Logger { return n }

func (n *noopLogger) WithName(_ string) contracts.Logger { return n }

func (n *noopLogger) WithOutput(_ io.Writer) contracts.Logger { return n }

func (n *noopLogger) WithPrefix(_ string) contracts.Logger { return n }

func (n *noopLogger) WithTracer(_ bool) contracts.Logger { return n }

func NewNoopLogger() contracts.Logger {
	return &noopLogger{}
}

func (n *noopLogger) Debug(...any)                    {}
func (n *noopLogger) Debugf(string, ...any)           {}
func (n *noopLogger) Debugw(string, contracts.Fields) {}
func (n *noopLogger) Info(...any)                     {}
func (n *noopLogger) Infof(string, ...any)            {}
func (n *noopLogger) Infow(string, contracts.Fields)  {}
func (n *noopLogger) Warn(...any)                     {}
func (n *noopLogger) Warnf(string, ...any)            {}
func (n *noopLogger) Warnw(string, contracts.Fields)  {}
func (n *noopLogger) Error(...any)                    {}
func (n *noopLogger) Errorf(string, ...any)           {}
func (n *noopLogger) Errorw(string, contracts.Fields) {}
func (n *noopLogger) Fatal(...any)                    {}
func (n *noopLogger) Fatalf(string, ...any)           {}
func (n *noopLogger) Fatalw(string, contracts.Fields) {}
func (n *noopLogger) Panic(...any)                    {}
func (n *noopLogger) Panicf(string, ...any)           {}
func (n *noopLogger) Panicw(string, contracts.Fields) {}
func (n *noopLogger) Print(...any)                    {}
func (n *noopLogger) Printf(string, ...any)           {}
func (n *noopLogger) Printw(string, contracts.Fields) {}

func (n *noopLogger) LogType() logtype.LogType {
	return logtype.LogTypes.NOOP
}

func (n *noopLogger) Instance() any {
	return nil
}
