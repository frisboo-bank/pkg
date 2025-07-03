package options

import (
	"time"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type RPCServerOption = func(options *RPCServerOptions)

func UseHost(host string) RPCServerOption {
	return func(o *RPCServerOptions) { o.Host = host }
}

func UsePort(port string) RPCServerOption {
	return func(o *RPCServerOptions) { o.Port = port }
}

func UseServerShutdownTimeout(timeout time.Duration) RPCServerOption {
	return func(o *RPCServerOptions) { o.ServerShutdownTimeout = timeout }
}

func UseServices(services Services) RPCServerOption {
	return func(o *RPCServerOptions) { o.Services = services }
}

func UseLogger(logger loggerContracts.Logger) RPCServerOption {
	return func(o *RPCServerOptions) { o.Logger = logger }
}
