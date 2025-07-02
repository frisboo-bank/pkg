package options

import (
	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/reflection/typemapper"
	"net"
	"strings"
	"time"

	"github.com/stoewer/go-strcase"
)

type RPCServerOptions struct {
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout"`
	Services              Services
	Logger                loggerContracts.Logger
}

type RPCServerOption = func(options *RPCServerOptions)

var DefaultRPCServerOptions = RPCServerOptions{
	Host: "0.0.0.0",
	Port: "9000",
}

func WithServices(services Services) RPCServerOption {
	return func(options *RPCServerOptions) {
		options.Services = services
	}
}

func WithLogger(logger loggerContracts.Logger) RPCServerOption {
	return func(options *RPCServerOptions) {
		options.Logger = logger
	}
}

func (o *RPCServerOptions) Address() string {
	return net.JoinHostPort(o.Host, o.Port)
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[RPCServerOptions]())

func ProvideRPCServerOptions(env environment.Environment) (*RPCServerOptions, error) {
	opt, err := config.BindConfigKey[*RPCServerOptions](optionName, env)
	if err != nil {
		return nil, err
	}

	return setDefaultOptions(opt), nil
}

func setDefaultOptions(options *RPCServerOptions) *RPCServerOptions {
	if strings.TrimSpace(options.Host) == "" {
		options.Host = DefaultRPCServerOptions.Host
	}

	if strings.TrimSpace(options.Port) == "" {
		options.Port = DefaultRPCServerOptions.Port
	}

	if options.ServerShutdownTimeout <= 0 {
		options.ServerShutdownTimeout = DefaultRPCServerOptions.ServerShutdownTimeout
	}

	return options
}
