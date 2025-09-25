package contracts

import (
	"context"

	"frisboo-bank/pkg/http/http_server/config"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	httpCommon interface {
		SetupDefaultMiddlewares()
		AddMiddlewares(middlewares ...any)
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
		ListRoutes() []any
		Name() string
		Type() httpservertype.HttpServerType
		Config() *config.Config
		RouteBuilder() RouteBuilder
		Logger() loggerContracts.Logger
	}

	HTTPServer interface {
		httpCommon
	}

	HTTPServerAdapter interface {
		httpCommon
	}
)
