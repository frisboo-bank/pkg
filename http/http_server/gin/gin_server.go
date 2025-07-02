package gin

import (
	"context"
	"errors"
	"fmt"
	"frisboo-bank/pkg/http/http_server/contracts"
	requestid "frisboo-bank/pkg/http/http_server/gin/middlewares/request_id"
	"frisboo-bank/pkg/http/http_server/options"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ginHttpServer struct {
	engine       *gin.Engine
	httpServer   *http.Server
	routeBuilder contracts.RouteBuilder
	logger       loggerContracts.Logger
	config       *options.HttpServerOptions
}

var _ contracts.HttpServer = (*ginHttpServer)(nil)

func NewGinHttpServer(config *options.HttpServerOptions) contracts.HttpServer {
	return newGinHttpServer(config)
}

func newGinHttpServer(config *options.HttpServerOptions) contracts.HttpServer {
	switch config.Development {
	case true:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()

	return &ginHttpServer{
		engine:       engine,
		routeBuilder: NewRouteBuilder(engine),
		logger:       config.Logger,
		config:       config,
	}
}

func (g *ginHttpServer) Start() error {
	g.httpServer = &http.Server{
		Addr:              g.config.Address(),
		Handler:           g.engine,
		ReadTimeout:       g.config.ReadTimeout,
		ReadHeaderTimeout: g.config.ReadHeaderTimeout,
		WriteTimeout:      g.config.WriteTimeout,
		IdleTimeout:       g.config.IdleTimeout,
		MaxHeaderBytes:    g.config.MaxHeaderBytes,
	}

	if err := g.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (g *ginHttpServer) Shutdown(ctx context.Context) error {
	if g.httpServer == nil {
		return fmt.Errorf("gin-server: looks like there is no server running")
	}

	return g.httpServer.Shutdown(ctx)
}

func (g *ginHttpServer) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(fmt.Errorf("gin-server: invalid middleware : `%v`", err))
	}

	g.engine.Use(ms...)
}

func (g *ginHttpServer) SetupDefaultMiddlewares() {
	g.AddMiddlewares(
		gin.Logger(),
		gin.Recovery(),
		cors.Default(),
		requestid.NewRequestIDMiddleware(),
	)
}

func (g *ginHttpServer) Instance() any {
	return g.engine
}

func (g *ginHttpServer) RouteBuilder() contracts.RouteBuilder {
	return g.routeBuilder
}

func (g *ginHttpServer) Config() *options.HttpServerOptions {
	return g.config
}

func (g *ginHttpServer) Logger() loggerContracts.Logger {
	return g.logger
}
