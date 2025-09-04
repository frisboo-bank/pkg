package config

import (
	"strings"
	"time"

	loggerConfig "frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/options"
	grpcConfig "frisboo-bank/pkg/rpc/rpc_server/adapters/grpc/config"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

type Option = options.OptionFn[Config]

var Type = options.OptionErr(func(c *Config, sType rpcservertype.RpcServerType) error {
	if err := validation.EnumOneOf("Type", sType, rpcservertype.RpcServerTypes); err != nil {
		return err
	}
	c.Type = sType
	return nil
})

var Host = options.OptionErr(func(c *Config, host string) error {
	host = strings.TrimSpace(host)
	if host == "" {
		return syserrors.CantBeEmptyError("Host")
	}
	c.Host = host
	return nil
})

var Port = options.OptionErr(func(c *Config, port string) error {
	port = strings.TrimSpace(port)
	if port == "" {
		return syserrors.CantBeEmptyError("Port")
	}
	c.Port = port
	return nil
})

var ServerShutdownTimeout = options.OptionErr(func(c *Config, serverShutdownTimeout time.Duration) error {
	if serverShutdownTimeout <= 0 {
		return syserrors.MustBePositiveError("ServerShutdownTimeout", serverShutdownTimeout)
	}
	c.ServerShutdownTimeout = serverShutdownTimeout
	return nil
})

var GRPC = options.OptionErr(func(c *Config, grpc *grpcConfig.Config) error {
	if grpc == nil {
		return syserrors.CantBeNilError("GRPC")
	}
	c.GRPC = grpc
	return nil
})

var Logger = options.OptionErr(func(c *Config, logger *loggerConfig.Config) error {
	if logger == nil {
		return syserrors.CantBeNilError("Logger")
	}
	c.Logger = logger
	return nil
})
