package contracts

import (
	"context"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type Container interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Logger() loggerContracts.Logger
}
