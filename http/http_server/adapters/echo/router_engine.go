package echo

import (
	"frisboo-bank/pkg/http/http_server/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	echoVendor "github.com/labstack/echo/v4"
)

var _ contracts.RouterEngine = (*echoRouterEngine)(nil)

type echoRouterEngine struct {
	echo   *echoVendor.Echo
	logger loggerContracts.Logger
}

func newRouterEngine(e *echoVendor.Echo, logger loggerContracts.Logger) contracts.RouterEngine {
	validation.AssertNotNil("echo", e)
	validation.AssertNotNil("logger", logger)

	return &echoRouterEngine{
		echo:   e,
		logger: logger,
	}
}

func (e *echoRouterEngine) Group(path string, middlewares ...any) contracts.RouterEngine {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid middleware for group:{%s}", path))
	}
	g := e.echo.Group(path, ms...)
	return newRouterEngine(g, e.logger)
}

func (e *echoRouterEngine) Handle(method string, path string, handler any, middlewares ...any) {
	h, err := ToHandlerFunc(handler)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid handler for route method:{%s} path:{%s}", method, path))
	}
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrapf(err, "invalid middleware for route method:{%s} path:{%s}", method, path))
	}
	e.echo.Add(method, path, h, ms...)
}

func (e *echoRouterEngine) Static(prefix string, root string) {
	e.echo.Static(prefix, root)
}
