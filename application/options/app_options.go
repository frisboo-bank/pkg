package options

import (
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
)

type AppOptions struct {
	Name             string `mapstructure:"name"`
	EnableHTTPServer bool   `mapstructure:"enableHTTPServer"`
	EnableGRPCServer bool   `mapstructure:"enableGRPCServer"`
}

func ProvideLoggerOptions(loader configContracts.ConfigLoader, env environment.Environment) (*AppOptions, error) {
	return config.LoadOptions[AppOptions](loader, env)
}
