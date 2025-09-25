package echo

import (
	"context"
	"strings"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	echoVendor "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"
)

var _ contracts.HTTPServerAdapter = (*echoHTTPServerAdapter)(nil)

type echoHTTPServerAdapter struct {
	cfg          *config.Config
	echo         *echoVendor.Echo
	logger       loggerContracts.Logger
	meter        metric.Meter
	routeBuilder contracts.RouteBuilder
}

func New(cfg *config.Config, logger loggerContracts.Logger, meter metric.Meter) contracts.HTTPServerAdapter {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)
	// validation.AssertNotNil("meter", meter)

	e := echoVendor.New()
	e.HideBanner = true

	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.ReadHeaderTimeout = cfg.ReadHeaderTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	e.Server.MaxHeaderBytes = cfg.MaxHeaderBytes

	return &echoHTTPServerAdapter{
		cfg:          cfg,
		echo:         e,
		logger:       logger,
		routeBuilder: NewRouteBuilder(e),
	}
}

func (e *echoHTTPServerAdapter) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrap(err, "invalid middleware"))
	}
	e.echo.Use(ms...)
}

func (e *echoHTTPServerAdapter) ListRoutes() []any {
	panic("unimplemented")
}

func (e *echoHTTPServerAdapter) SetupDefaultMiddlewares() {
	// TODO: improve to support pattern matching
	skipper := func(c echoVendor.Context) bool {
		rPath := c.Request().URL.Path
		for _, skip := range e.cfg.IgnoreLogUrls {
			if strings.Contains(rPath, skip) {
				return true
			}
		}
		return false
	}

	e.echo.Use(
		middleware.Recover(),
		middleware.BodyLimit(e.cfg.BodyLimit),
		middleware.RequestID(),
		// middlewares.IPRateLimit(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: skipper,
			Level:   e.cfg.GzipLevel,
		}),
	)
}

func (e *echoHTTPServerAdapter) Start(ctx context.Context) error {
	return e.echo.Start(e.cfg.Address())
}

func (e *echoHTTPServerAdapter) Stop(ctx context.Context) error {
	return e.echo.Shutdown(ctx)
}

func (e *echoHTTPServerAdapter) Type() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.ECHO
}

func (e *echoHTTPServerAdapter) Config() *config.Config {
	return e.cfg
}

func (e *echoHTTPServerAdapter) Logger() loggerContracts.Logger {
	return e.logger
}

func (e *echoHTTPServerAdapter) RouteBuilder() contracts.RouteBuilder {
	return e.routeBuilder
}
