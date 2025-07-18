package contracts

import (
	"context"
	"time"

	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type HTTPServer interface {
	WithBasePath(base string) HTTPServer
	HasDevelopment(dev bool) HTTPServer
	WithHost(host string) HTTPServer
	WithPort(port string) HTTPServer
	WithIgnoreLogUrls(urls []string) HTTPServer
	WithBodyLimit(limit string) HTTPServer
	WithIdleTimeout(timeout time.Duration) HTTPServer
	WithMaxHeaderBytes(max int) HTTPServer
	WithReadHeaderTimeout(timeout time.Duration) HTTPServer
	WithReadTimeout(timeout time.Duration) HTTPServer
	WithServerShutdownTimeout(timeout time.Duration) HTTPServer
	WithWriteTimeout(timeout time.Duration) HTTPServer

	Start() error
	Shutdown(ctx context.Context) error
	AddMiddlewares(middlewares ...any)
	SetupDefaultMiddlewares()
	Instance() any
	ServerType() httpservertype.HTTPServerType
	Address() string

	RouteBuilder() RouteBuilder
	Logger() loggerContracts.Logger
}
