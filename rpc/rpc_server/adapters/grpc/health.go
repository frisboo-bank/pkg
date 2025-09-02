package grpc

import (
	"context"

	"frisboo-bank/pkg/rpc/rpc_server/contracts"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	healthContracts "frisboo-bank/pkg/health/contracts"
)

var _ healthContracts.HealthServiceCheck = (*GRPCHealthService)(nil)

type GRPCHealthService struct {
	client contracts.RPCServer
	health string
	logger loggerContracts.Logger
}

func NewGRPCHealthService(client contracts.RPCServer) healthContracts.HealthServiceCheck {
	return &GRPCHealthService{
		client: client,
	}
}

func (g *GRPCHealthService) CheckHealth(ctx context.Context) error {
	return nil
}

func (g *GRPCHealthService) GetServiceName() string {
	return "grpc"
}

func (g *GRPCHealthService) Logger() loggerContracts.Logger {
	return g.logger
}
