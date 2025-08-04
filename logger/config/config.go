package config

import (
	"os"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"

	configContracts "frisboo-bank/pkg/config/contracts"

	encodingtype "frisboo-bank/pkg/logger/contracts/enums/encoding_type"
	loglevel "frisboo-bank/pkg/logger/contracts/enums/log_level"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

var (
	Level    = loglevel.LogLevels.ERRORLEVEL
	Encoding = encodingtype.EncodingTypes.TEXT
	Output   = os.Stdout
)

type LoggerConfig struct {
	Type          loggertype.LoggerType     `mapstructure:"type"`
	Level         loglevel.LogLevel         `mapstructure:"level"`
	CallerEnabled bool                      `mapstructure:"callerEnabled"`
	EnableTracing bool                      `mapstructure:"enableTracing"`
	CallDepth     int                       `mapstructure:"callDepth"`
	Encoding      encodingtype.EncodingType `mapstructure:"encoding"`
}

func ProvideLoggerConfig(loader configContracts.ConfigLoader, env environment.Environment) (*LoggerConfig, error) {
	return config.LoadConfig[LoggerConfig](loader, env, "logger")
}
