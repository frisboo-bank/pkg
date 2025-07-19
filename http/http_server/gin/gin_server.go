package gin

import (
	"context"
	"errors"
	"fmt"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/options"
	"net"
	"net/http"
	"time"

	requestid "frisboo-bank/pkg/http/http_server/gin/middlewares/request_id"

	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ginHTTPServer struct {
	basePath              string
	bodyLimit             string
	development           bool
	host                  string
	idleTimeout           time.Duration
	ignoreLogUrls         []string
	logger                loggerContracts.Logger
	maxHeaderBytes        int
	port                  string
	readHeaderTimeout     time.Duration
	readTimeout           time.Duration
	serverShutdownTimeout time.Duration
	writeTimeout          time.Duration

	engine       *gin.Engine
	httpServer   *http.Server
	routeBuilder contracts.RouteBuilder
}

func (g *ginHTTPServer) WithOptions(options *options.HTTPServerOptions) contracts.HTTPServer {
	return g.
		WithBasePath(options.BasePath).
		WithBodyLimit(options.BodyLimit).
		HasDevelopment(options.Development).
		WithHost(options.Host).
		WithIdleTimeout(options.IdleTimeout).
		WithIgnoreLogUrls(options.IgnoreLogUrls).
		WithMaxHeaderBytes(options.MaxHeaderBytes).
		WithPort(options.Port).
		WithReadHeaderTimeout(options.ReadHeaderTimeout).
		WithReadTimeout(options.ReadTimeout).
		WithServerShutdownTimeout(options.ServerShutdownTimeout).
		WithWriteTimeout(options.WriteTimeout)
}

func (g *ginHTTPServer) WithBasePath(base string) contracts.HTTPServer {
	g.basePath = base
	return g
}

func (g *ginHTTPServer) WithBodyLimit(limit string) contracts.HTTPServer {
	g.bodyLimit = limit
	return g
}

func (g *ginHTTPServer) HasDevelopment(dev bool) contracts.HTTPServer {
	g.development = dev
	return g
}

func (g *ginHTTPServer) WithHost(host string) contracts.HTTPServer {
	g.host = host
	return g
}

func (g *ginHTTPServer) WithIdleTimeout(timeout time.Duration) contracts.HTTPServer {
	g.idleTimeout = timeout
	return g
}

func (g *ginHTTPServer) WithIgnoreLogUrls(urls []string) contracts.HTTPServer {
	g.ignoreLogUrls = urls
	return g
}

func (g *ginHTTPServer) WithMaxHeaderBytes(max int) contracts.HTTPServer {
	g.maxHeaderBytes = max
	return g
}

func (g *ginHTTPServer) WithPort(port string) contracts.HTTPServer {
	g.port = port
	return g
}

func (g *ginHTTPServer) WithReadHeaderTimeout(timeout time.Duration) contracts.HTTPServer {
	g.readHeaderTimeout = timeout
	return g
}

func (g *ginHTTPServer) WithReadTimeout(timeout time.Duration) contracts.HTTPServer {
	g.readTimeout = timeout
	return g
}

func (g *ginHTTPServer) WithServerShutdownTimeout(timeout time.Duration) contracts.HTTPServer {
	g.serverShutdownTimeout = timeout
	return g
}

func (g *ginHTTPServer) WithWriteTimeout(timeout time.Duration) contracts.HTTPServer {
	g.writeTimeout = timeout
	return g
}

var _ contracts.HTTPServer = (*ginHTTPServer)(nil)

func NewGinHTTPServer(logger loggerContracts.Logger) contracts.HTTPServer {
	engine := gin.New()

	return &ginHTTPServer{
		basePath:              options.BasePath,
		bodyLimit:             options.BodyLimit,
		development:           false,
		host:                  options.Host,
		idleTimeout:           options.IdleTimeout,
		logger:                logger,
		maxHeaderBytes:        options.MaxHeaderBytes,
		port:                  options.Port,
		readHeaderTimeout:     options.ReadHeaderTimeout,
		readTimeout:           options.ReadTimeout,
		serverShutdownTimeout: options.ServerShutdownTimeout,
		writeTimeout:          options.WriteTimeout,

		engine:       engine,
		routeBuilder: NewRouteBuilder(engine),
	}
}

func (g *ginHTTPServer) Start() error {
	switch g.development {
	case true:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	g.httpServer = &http.Server{
		Addr:              g.Address(),
		Handler:           g.engine,
		ReadTimeout:       g.readTimeout,
		ReadHeaderTimeout: g.readHeaderTimeout,
		WriteTimeout:      g.writeTimeout,
		IdleTimeout:       g.idleTimeout,
		MaxHeaderBytes:    g.maxHeaderBytes,
	}

	if err := g.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (g *ginHTTPServer) Shutdown(ctx context.Context) error {
	if g.httpServer == nil {
		return fmt.Errorf("gin-server: looks like there is no server running")
	}

	return g.httpServer.Shutdown(ctx)
}

func (g *ginHTTPServer) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(fmt.Errorf("gin-server: invalid middleware : `%v`", err))
	}

	g.engine.Use(ms...)
}

func (g *ginHTTPServer) SetupDefaultMiddlewares() {
	g.AddMiddlewares(
		gin.Logger(),
		gin.Recovery(),
		cors.Default(),
		requestid.NewRequestIDMiddleware(),
	)
}

func (g *ginHTTPServer) ServerType() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.GIN
}

func (g *ginHTTPServer) Instance() any {
	return g.engine
}

func (g *ginHTTPServer) Address() string {
	return net.JoinHostPort(g.host, g.port)
}

func (g *ginHTTPServer) RouteBuilder() contracts.RouteBuilder {
	return g.routeBuilder
}

func (g *ginHTTPServer) Logger() loggerContracts.Logger {
	return g.logger
}
