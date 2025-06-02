package contracts

import (
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/waiter"
)

type Application interface {
	Di() di.Container
	Logger() logger.Logger
	Environment() environment.Environment
	Waiter() waiter.Waiter
	Run() error
}
