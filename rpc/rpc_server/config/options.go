package config

import (
	"strings"
	"time"

	"frisboo-bank/pkg/options"
	grpcConfig "frisboo-bank/pkg/rpc/rpc_server/adapters/grpc/config"
	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/enums/rpc_server_type"
)

type Option = options.OptionFn[Config]

var Type = options.Option(func(c *Config, sType rpcservertype.RpcServerType) {
	c.Type = sType
})

var Host = options.Option(func(c *Config, host string) {
	c.Host = strings.TrimSpace(host)
})

var Port = options.Option(func(c *Config, port string) {
	c.Port = strings.TrimSpace(port)
})

var ServerShutdownTimeout = options.Option(func(c *Config, serverShutdownTimeout time.Duration) {
	c.ServerShutdownTimeout = serverShutdownTimeout
})

var GRPC = options.Option(func(c *Config, grpc *grpcConfig.Config) {
	c.GRPC = grpc
})

var Logger = options.Option(func(c *Config, logger string) {
	c.Logger = logger
})
