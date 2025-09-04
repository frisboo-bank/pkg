package contracts

import (
	"context"

	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	httpCommon interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		Type() httpservertype.HttpServerType
		RouteBuilder() RouteBuilder
		Logger() loggerContracts.Logger
		ListRoutes() []any
	}

	HTTPServer interface {
		httpCommon
	}

	HTTPServerAdapter interface {
		httpCommon
	}
)
