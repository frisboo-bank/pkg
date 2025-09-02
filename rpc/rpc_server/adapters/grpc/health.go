package grpc

import (
	"context"

	healthContracts "frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/rpc/rpc_server/contracts"
)

type GRPCHealthService struct {
	client contracts.RPCServer
	health string
}

var _ healthContracts.HealthServiceCheck = (*GRPCHealthService)(nil)

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
