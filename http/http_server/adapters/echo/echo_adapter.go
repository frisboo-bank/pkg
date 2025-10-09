package echo

import (
	"context"
	"fmt"
	"strings"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	"frisboo-bank/pkg/http/http_server/routing"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	echoVendor "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel/metric"
)

var _ contracts.HTTPServerAdapter = (*echoHTTPServerAdapter)(nil)

type echoHTTPServerAdapter struct {
	name         string
	cfg          *config.Config
	echo         *echoVendor.Echo
	logger       loggerContracts.Logger
	meter        metric.Meter
	routeBuilder contracts.RouteBuilder
}

func New(
	name string,
	cfg *config.Config,
	logger loggerContracts.Logger,
	meter metric.Meter,
) contracts.HTTPServerAdapter {
	validation.AssertNotEmpty("name", name)
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

	routerEngine := newRouterEngine(e, logger)

	return &echoHTTPServerAdapter{
		name:         name,
		cfg:          cfg,
		echo:         e,
		logger:       logger,
		meter:        meter,
		routeBuilder: routing.NewBuilder(routerEngine),
	}
}

func (e *echoHTTPServerAdapter) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrap(err, "invalid middleware"))
	}
	e.echo.Use(ms...)
}

func (e *echoHTTPServerAdapter) SetupDefaultMiddlewares() {
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
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			Skipper: skipper,

			LogRequestID:     true,
			LogRemoteIP:      true,
			LogHost:          true,
			LogMethod:        true,
			LogURI:           true,
			LogUserAgent:     true,
			LogStatus:        true,
			LogError:         true,
			LogLatency:       true,
			LogContentLength: true,
			LogResponseSize:  true,

			LogValuesFunc: func(c echoVendor.Context, v middleware.RequestLoggerValues) error {
				e.logger.Infow(fmt.Sprintf("[Request] url:{%v} status:{%v}", v.URI, v.Status),
					loggerContracts.Fields{
						"uri":           v.URI,
						"status":        v.Status,
						"id":            v.RequestID,
						"remote_ip":     v.RemoteIP,
						"host":          v.Host,
						"method":        v.Method,
						"user_agent":    v.UserAgent,
						"error":         v.Error,
						"latency":       v.Latency.Nanoseconds(),
						"latency_human": v.Latency.String(),
						"bytes_in":      v.ContentLength,
						"bytes_out":     v.ResponseSize,
					})

				return nil
			},
		}),
		middleware.Recover(),
		middleware.AddTrailingSlash(),
		middleware.BodyLimit(e.cfg.BodyLimit),
		middleware.RequestID(),
		middleware.Secure(),
		middleware.GzipWithConfig(middleware.GzipConfig{
			Skipper: skipper,
			Level:   e.cfg.GzipLevel,
		}),
		middleware.CORS(),
	)
}

func (e *echoHTTPServerAdapter) Start(ctx context.Context) error {
	return e.echo.Start(e.cfg.Address())
}

func (e *echoHTTPServerAdapter) Stop(ctx context.Context) error {
	return e.echo.Shutdown(ctx)
}

func (e *echoHTTPServerAdapter) ListRoutes() []any {
	rs := e.echo.Routes()
	out := make([]any, len(rs))
	for i, r := range rs {
		out[i] = r
	}
	return out
}

func (e *echoHTTPServerAdapter) Name() string {
	return e.name
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
