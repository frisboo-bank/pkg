package contracts

import (
	"net"
	"time"

	"frisboo-bank/pkg/http/http_server/config"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type BaseHTTPServer struct {
	basePath              string
	bodyLimit             string
	development           bool
	host                  string
	idleTimeout           time.Duration
	ignoreLogUrls         []string
	internal              HTTPServerInternal
	logger                loggerContracts.Logger
	maxHeaderBytes        int
	middlewares           []any
	port                  string
	readHeaderTimeout     time.Duration
	readTimeout           time.Duration
	routeBuilder          RouteBuilder
	serverShutdownTimeout time.Duration
	writeTimeout          time.Duration
}

var _ httpServerConfig = (*BaseHTTPServer)(nil)

func (b *BaseHTTPServer) Init(internal HTTPServerInternal) {
	b.internal = internal
}

func (b *BaseHTTPServer) Address() string                      { return net.JoinHostPort(b.host, b.port) }
func (b *BaseHTTPServer) BasePath() string                     { return b.basePath }
func (b *BaseHTTPServer) BodyLimit() string                    { return b.bodyLimit }
func (b *BaseHTTPServer) Development() bool                    { return b.development }
func (b *BaseHTTPServer) Host() string                         { return b.host }
func (b *BaseHTTPServer) IdleTimeout() time.Duration           { return b.idleTimeout }
func (b *BaseHTTPServer) IgnoreLogUrls() []string              { return b.ignoreLogUrls }
func (b *BaseHTTPServer) Logger() loggerContracts.Logger       { return b.logger }
func (b *BaseHTTPServer) MaxHeaderBytes() int                  { return b.maxHeaderBytes }
func (b *BaseHTTPServer) Middlewares() []any                   { return b.middlewares }
func (b *BaseHTTPServer) Port() string                         { return b.port }
func (b *BaseHTTPServer) ReadHeaderTimeout() time.Duration     { return b.readHeaderTimeout }
func (b *BaseHTTPServer) ReadTimeout() time.Duration           { return b.readTimeout }
func (b *BaseHTTPServer) RouteBuilder() RouteBuilder           { return b.routeBuilder }
func (b *BaseHTTPServer) ServerShutdownTimeout() time.Duration { return b.serverShutdownTimeout }
func (b *BaseHTTPServer) WriteTimeout() time.Duration          { return b.writeTimeout }

func (b *BaseHTTPServer) WithConfig(cfg *config.HTTPServerConfig) HTTPServer {
	b.basePath = cfg.BasePath
	b.bodyLimit = cfg.BodyLimit
	b.development = cfg.Development
	b.host = cfg.Host
	b.idleTimeout = cfg.IdleTimeout
	b.ignoreLogUrls = cfg.IgnoreLogUrls
	b.maxHeaderBytes = cfg.MaxHeaderBytes
	b.port = cfg.Port
	b.readHeaderTimeout = cfg.ReadHeaderTimeout
	b.readTimeout = cfg.ReadTimeout
	b.serverShutdownTimeout = cfg.ServerShutdownTimeout
	b.writeTimeout = cfg.WriteTimeout

	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithBasePath(basePath string) HTTPServer {
	b.basePath = basePath
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithBodyLimit(bodyLimit string) HTTPServer {
	b.bodyLimit = bodyLimit
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithDevelopmentEnv(development bool) HTTPServer {
	b.development = development
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithHost(host string) HTTPServer {
	b.host = host
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithIdleTimeout(idleTimeout time.Duration) HTTPServer {
	b.idleTimeout = idleTimeout
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithIgnoreLogUrls(ignoreLogUrls []string) HTTPServer {
	b.ignoreLogUrls = ignoreLogUrls
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithLogger(logger loggerContracts.Logger) HTTPServer {
	b.logger = logger
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithMaxHeaderBytes(maxHeaderBytes int) HTTPServer {
	b.maxHeaderBytes = maxHeaderBytes
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithMiddlewares(middlewares ...any) HTTPServer {
	b.middlewares = middlewares
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithPort(port string) HTTPServer {
	b.port = port
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithReadHeaderTimeout(readHeaderTimeout time.Duration) HTTPServer {
	b.readHeaderTimeout = readHeaderTimeout
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithReadTimeout(readTimeout time.Duration) HTTPServer {
	b.readTimeout = readTimeout
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithRouteBuilder(routeBuilder RouteBuilder) HTTPServer {
	b.routeBuilder = routeBuilder
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithServerShutdownTimeout(serverShutdownTimeout time.Duration) HTTPServer {
	b.serverShutdownTimeout = serverShutdownTimeout
	return b.internal.(HTTPServer)
}

func (b *BaseHTTPServer) WithWriteTimeout(writeTimeout time.Duration) HTTPServer {
	b.writeTimeout = writeTimeout
	return b.internal.(HTTPServer)
}
