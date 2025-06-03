package contracts

import (
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/logger"
)

type System interface {
	Di() di.Container
	Logger() logger.Logger
	Environment() environment.Environment
	Run() error
}
