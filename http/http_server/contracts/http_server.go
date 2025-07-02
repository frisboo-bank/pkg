package contracts

import (
	"context"

	"frisboo-bank/pkg/http/http_server/options"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type HttpServer interface {
	Start() error
	Shutdown(ctx context.Context) error
	AddMiddlewares(middlewares ...any)
	SetupDefaultMiddlewares()
	Instance() any
	RouteBuilder() RouteBuilder
	Logger() loggerContracts.Logger
	Config() *options.HttpServerOptions
}
