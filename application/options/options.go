package options

import (
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
)

type AppConfig struct {
	Name             string `mapstructure:"name"`
	EnableHTTPServer bool   `mapstructure:"enableHTTPServer"`
	EnableGRPCServer bool   `mapstructure:"enableGRPCServer"`
}

func ProvideAppOptions(loader configContracts.ConfigLoader, env environment.Environment) (*AppConfig, error) {
	return config.LoadConfig[AppConfig](loader, env, "app")
}
