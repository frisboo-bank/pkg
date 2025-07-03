package options

import (
	"fmt"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"

	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"

	optionsContracts "frisboo-bank/pkg/options/contracts"
)

// LogOptions holds the logger configuration.
type LogOptions struct {
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Type          logtype.LogType           `mapstructure:"type"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	EnableTracing bool                      `mapstructure:"enableTracing"`
	CallDepth     int                       `mapstructure:"callDepth"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
}

var defaultLogOptions = &LogOptions{
	Level:    loglevel.LogLevels.ERROR_LEVEL,
	Type:     logtype.LogTypes.NOOP,
	Encoding: encodingtype.EncodingTypes.TEXT,
}

var _ optionsContracts.Options = (*LogOptions)(nil)

func (o *LogOptions) Clone() optionsContracts.Options {
	return options.CloneOptions(o)
}

func (o *LogOptions) SetDefaults() {
	if !o.Level.IsValid() {
		o.Level = defaultLogOptions.Level
	}

	if !o.Type.IsValid() {
		o.Type = defaultLogOptions.Type
	}

	if !o.Encoding.IsValid() {
		o.Encoding = defaultLogOptions.Encoding
	}
}

func (o *LogOptions) Validate() error {
	if !o.Level.IsValid() {
		return fmt.Errorf("(log-options) invalid log level: %q", o.Level)
	}
	if !o.Type.IsValid() {
		return fmt.Errorf("(log-options) invalid log type: %q", o.Type)
	}
	if !o.Encoding.IsValid() {
		return fmt.Errorf("(log-options) invalid encoding type: %q", o.Encoding)
	}

	return nil
}

func ProvideLogOptions(env environment.Environment) (*LogOptions, error) {
	return options.LoadOptions[*LogOptions](env, func(options *LogOptions) { options.SetDefaults() })
}
