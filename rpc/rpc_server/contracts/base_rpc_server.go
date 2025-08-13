package contracts

import (
	"net"
	"time"

	"frisboo-bank/pkg/rpc/rpc_server/config"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type BaseRPCServer struct {
	host                  string
	port                  string
	serverShutdownTimeout time.Duration
	services              []Services
	logger                loggerContracts.Logger

	internal RPCServerInternal
}

var _ rpcServerConfig = (*BaseRPCServer)(nil)

func (b *BaseRPCServer) Init(internal RPCServerInternal) {
	b.internal = internal
}

func (b *BaseRPCServer) WithConfig(cfg *config.RPCServerConfig) RPCServer {
	b.host = cfg.Host
	b.port = cfg.Port
	b.serverShutdownTimeout = cfg.ServerShutdownTimeout

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) WithHost(host string) RPCServer {
	b.host = host

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) WithPort(port string) RPCServer {
	b.port = port

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) WithServerShutdownTimeout(serverShutdownTimeout time.Duration) RPCServer {
	b.serverShutdownTimeout = serverShutdownTimeout

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) WithServices(services []Services) RPCServer {
	b.services = append(b.services, services...)

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) WithLogger(logger loggerContracts.Logger) RPCServer {
	b.logger = logger

	b.internal.SetupInstance()
	return b.internal.(RPCServer)
}

func (b *BaseRPCServer) Host() string                         { return b.host }
func (b *BaseRPCServer) Port() string                         { return b.port }
func (b *BaseRPCServer) ServerShutdownTimeout() time.Duration { return b.serverShutdownTimeout }
func (b *BaseRPCServer) Services() []Services                 { return b.services }
func (b *BaseRPCServer) Logger() loggerContracts.Logger       { return b.logger }
func (b *BaseRPCServer) Address() string                      { return net.JoinHostPort(b.host, b.port) }
