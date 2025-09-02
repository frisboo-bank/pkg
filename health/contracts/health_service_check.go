package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type HealthServiceCheck interface {
	CheckHealth(ctx context.Context) error
	GetServiceName() string
	Logger() loggerContracts.Logger
}
