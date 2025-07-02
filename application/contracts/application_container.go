package contracts

import (
	"context"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type ApplicationContainer interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Logger() loggerContracts.Logger
}
