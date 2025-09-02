package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type HealthService interface {
	AddServices(services ...HealthServiceCheck)
	CheckHealth(ctx context.Context) CheckAllStatus
	Logger() loggerContracts.Logger
}
