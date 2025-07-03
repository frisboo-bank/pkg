package options

import (
	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

// LogOption is a functional option for LogOptions.
type LogOption func(options *LogOptions)

// NewLogOptions constructs LogOptions with safe defaults and applies any LogOption.
func NewLogOptions(options ...LogOption) *LogOptions {
	config := *defaultLogOptions

	for _, opt := range options {
		opt(&config)
	}

	return &config
}

// WithLevel sets the log level.
func WithLevel(level loglevel.LogLevel) LogOption {
	return func(o *LogOptions) { o.Level = level }
}

// WithType sets the logger backend type.
func WithType(t logtype.LogType) LogOption {
	return func(o *LogOptions) { o.Type = t }
}

// WithCaller enables or disables caller reporting and sets call stack depth.
func WithCaller(enabled bool, depth int) LogOption {
	return func(o *LogOptions) {
		o.CallerEnabled = enabled
		o.CallDepth = depth
	}
}

// WithTracer enables or disables tracing.
func WithTracer(enabled bool) LogOption {
	return func(o *LogOptions) { o.EnableTracing = enabled }
}

// WithEncoding sets the output encoding.
func WithEncoding(encoding encodingtype.EncodingType) LogOption {
	return func(o *LogOptions) { o.Encoding = encoding }
}
