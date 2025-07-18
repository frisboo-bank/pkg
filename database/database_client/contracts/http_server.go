package contracts

import (
	"context"

	"frisboo-bank/pkg/http/http_server/config"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type HTTPServer interface {
	Start() error
	Shutdown(ctx context.Context) error
	AddMiddlewares(middlewares ...any)
	SetupDefaultMiddlewares()
	Instance() any
	RouteBuilder() RouteBuilder
	Logger() loggerContracts.Logger
	Config() *config.HTTPServerOptions
}
