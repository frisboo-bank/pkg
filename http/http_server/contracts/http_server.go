package contracts

import (
	"context"
	"time"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/http/http_server/config"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
)

var (
	Type                  = httpservertype.HttpServerTypes.GIN
	Host                  = "0.0.0.0"
	Port                  = "8000"
	ServerShutdownTimeout = constants.SERVER_SHUTDOWN_TIMEOUT
)

type (
	httpServerConfig interface {
		BasePath() string
		BodyLimit() string
		Development() bool
		Host() string
		IdleTimeout() time.Duration
		IgnoreLogUrls() []string
		Logger() loggerContracts.Logger
		MaxHeaderBytes() int
		Middlewares() []any
		Port() string
		ReadHeaderTimeout() time.Duration
		ReadTimeout() time.Duration
		RouteBuilder() RouteBuilder
		ServerShutdownTimeout() time.Duration
		WriteTimeout() time.Duration

		WithBasePath(basePath string) HTTPServer
		WithBodyLimit(bodyLimit string) HTTPServer
		WithConfig(config *config.HTTPServerConfig) HTTPServer
		WithDevelopmentEnv(dev bool) HTTPServer
		WithHost(host string) HTTPServer
		WithIdleTimeout(idleTimeout time.Duration) HTTPServer
		WithIgnoreLogUrls(ignoreLogUrls []string) HTTPServer
		WithLogger(logger loggerContracts.Logger) HTTPServer
		WithMaxHeaderBytes(maxHeaderBytes int) HTTPServer
		WithMiddlewares(middlewares ...any) HTTPServer
		WithPort(port string) HTTPServer
		WithReadHeaderTimeout(readHeaderTimeout time.Duration) HTTPServer
		WithReadTimeout(readTimeout time.Duration) HTTPServer
		WithRouteBuilder(routeBuilder RouteBuilder) HTTPServer
		WithServerShutdownTimeout(serverShutdownTimeout time.Duration) HTTPServer
		WithWriteTimeout(writeTimeout time.Duration) HTTPServer
	}

	httpServerCore interface {
		Address() string
		Instance() any
		SetupDefaultMiddlewares()
		Shutdown(ctx context.Context) error
		Start() error
		Type() httpservertype.HttpServerType
	}

	HTTPServer interface {
		httpServerConfig
		httpServerCore
	}

	HTTPServerInternal interface {
		SetupInstance()
	}
)
