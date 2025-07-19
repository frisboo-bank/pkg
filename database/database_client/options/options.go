package options

import (
	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
)

type (
	Type string
)

const (
	TypePostgres = Type("postgres")
)

type DatabaseClientOptions struct {
	Type     Type   `mapstructure:"type"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	SSLMode  bool   `mapstructure:"sslMode"`
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

func ProvideDatabaseClientOptions(
	loader configContracts.ConfigLoader,
	env environment.Environment,
) (*DatabaseClientOptions, error) {
	return config.LoadOptions[DatabaseClientOptions](loader, env)
}
