package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"frisboo-bank/pkg/rpc/rpc_server/contracts"

	rpcservertype "frisboo-bank/pkg/rpc/rpc_server/contracts/enums/rpc_server_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	googlerpc "google.golang.org/grpc"
)

type GRPCServer struct {
	contracts.BaseRPCServer
	instance *googlerpc.Server
}

func New(logger loggerContracts.Logger) contracts.RPCServer {
	rpcServer := &GRPCServer{
		instance: googlerpc.NewServer(),
	}
	rpcServer.Init(rpcServer)
	rpcServer.SetupInstance()
	rpcServer.WithLogger(logger)
	return rpcServer
}

func (s *GRPCServer) SetupInstance() {
}

func (s *GRPCServer) Start(listener net.Listener) error {
	if err := s.instance.Serve(listener); err != nil && !errors.Is(err, googlerpc.ErrServerStopped) {
		return err
	}

	return nil
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	if s.instance == nil {
		return fmt.Errorf("grpc-server: looks like there is no server running")
	}

	s.instance.GracefulStop()

	return nil
}

func (s *GRPCServer) Instance() any {
	return s.instance
}

func (s *GRPCServer) Type() rpcservertype.RpcServerType {
	return rpcservertype.RpcServerTypes.GRPC
}
