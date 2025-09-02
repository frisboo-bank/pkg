package config

import (
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	databaseclienttype "frisboo-bank/pkg/database/database_client/contracts/enums/database_client_type"
	"frisboo-bank/pkg/environment"
)

type dbEnvConfig struct {
	Type     databaseclienttype.DatabaseClientType `mapstructure:"type"`
	Host     string                                `mapstructure:"host"`
	Port     string                                `mapstructure:"port"`
	User     string                                `mapstructure:"user"`
	Password string                                `mapstructure:"password"`
	SSLMode  bool                                  `mapstructure:"sslMode"`
}

type EnvConfig struct {
	DB map[string]dbEnvConfig `mapstructure:"db"`
}

func LoadEnvConfig(loader configloaderContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.Load[EnvConfig](loader, env, "db")
}

// type DatabaseClientOption = func(options *DatabaseClientOptions)
//
// var DefaultDatabaseClientOptions = DatabaseClientOptions{
// 	Host:    "127.0.0.1",
// 	SSLMode: false,
// 	Logger:  noop.NewNoopLogger(),
// }
//
// func WithOptions(partialOptions *DatabaseClientOptions) DatabaseClientOption {
// 	return func(options *DatabaseClientOptions) {
// 		if partialOptions.Type != "" {
// 			options.Type = partialOptions.Type
// 		}
//
// 		if partialOptions.Host != "" {
// 			options.Host = partialOptions.Host
// 		}
//
// 		if partialOptions.Port != "" {
// 			options.Port = partialOptions.Port
// 		}
//
// 		if partialOptions.User != "" {
// 			options.User = partialOptions.User
// 		}
//
// 		if partialOptions.Password != "" {
// 			options.Password = partialOptions.Password
// 		}
//
// 		if partialOptions.SSLMode {
// 			options.SSLMode = partialOptions.SSLMode
// 		}
// 	}
// }
//
// func WithLogger(logger loggerContracts.Logger) DatabaseClientOption {
// 	return func(options *DatabaseClientOptions) {
// 		options.Logger = logger
// 	}
// }
