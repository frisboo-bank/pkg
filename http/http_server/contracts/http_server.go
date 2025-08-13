package contracts

import (
	"context"

	"frisboo-bank/pkg/http/http_server/config"
	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	HTTPServer interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Shutdown(ctx context.Context) error
		Type() httpservertype.HttpServerType
		RouteBuilder() RouteBuilder
		Logger() loggerContracts.Logger
	}

	HTTPServerAdapter interface {
		Setup(cfg *config.Config) error
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Shutdown(ctx context.Context) error
		RouteBuilder() RouteBuilder
		Type() httpservertype.HttpServerType
	}
)

type (
	httpServerCore interface {
		SetupDefaultMiddlewares()
	}
)
