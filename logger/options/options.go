package options

import (
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
	"os"

	encodingtype "frisboo-bank/pkg/logger/options/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/options/enums/log_level"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

var (
	Level    = loglevel.LogLevels.ERROR_LEVEL
	Encoding = encodingtype.EncodingTypes.TEXT
	Output   = os.Stdout
)

type LogOptions struct {
	Level         loglevel.LogLevel         `mapstructure:"level"`
	Type          logtype.LogType           `mapstructure:"type"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	EnableTracing bool                      `mapstructure:"enableTracing"`
	CallDepth     int                       `mapstructure:"callDepth"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
}

func ProvideLogOptions(loader configContracts.ConfigLoader, env environment.Environment) (*LogOptions, error) {
	return config.LoadOptions[LogOptions](loader, env)
}
