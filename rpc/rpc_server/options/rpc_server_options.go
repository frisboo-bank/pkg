package options

import (
	"fmt"
	"net"
	"time"

	"frisboo-bank/pkg/environment"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	optionsContracts "frisboo-bank/pkg/options/contracts"
	"frisboo-bank/pkg/reflection/typemapper"

	"github.com/stoewer/go-strcase"
)

type RPCServerOptions struct {
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	ServerShutdownTimeout time.Duration `mapstructure:"serverShutdownTimeout"`
	Services              Services
	Logger                loggerContracts.Logger
}

var defaultOptions = RPCServerOptions{
	Host: "0.0.0.0",
	Port: "9000",
}

func (o *RPCServerOptions) Clone() optionsContracts.Options {
	return options.CloneOptions(o)
}

func (o *RPCServerOptions) SetDefaults() {
	if o.Host == "" {
		o.Host = defaultOptions.Host
	}

	if o.Port == "" {
		o.Port = defaultOptions.Port
	}

	if o.ServerShutdownTimeout <= 0 {
		o.ServerShutdownTimeout = defaultOptions.ServerShutdownTimeout
	}
}

func (o *RPCServerOptions) Validate() error {
	if o.Host == "" {
		return fmt.Errorf("Host must not be empty")
	}

	if o.Port == "" {
		return fmt.Errorf("Port must not be empty")
	}
	return nil
}

func (o *RPCServerOptions) Address() string {
	return net.JoinHostPort(o.Host, o.Port)
}

var optionName = strcase.LowerCamelCase(typemapper.GetGenericTypeNameByT[RPCServerOptions]())

// ApplyRPCServerOptions applies a series of RPCServerOption functions to an RPCServerOptions struct.
func ApplyRPCServerOptions(o *RPCServerOptions, opts ...RPCServerOption) {
	for _, opt := range opts {
		opt(o)
	}
}

func ProvideRPCServerOptions(env environment.Environment) (*RPCServerOptions, error) {
	return options.LoadOptions[*RPCServerOptions](env, func(options *RPCServerOptions) { options.SetDefaults() })
}
