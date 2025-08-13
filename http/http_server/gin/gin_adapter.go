package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	requestid "frisboo-bank/pkg/http/http_server/gin/middlewares/request_id"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginVendor "github.com/gin-gonic/gin"
)

var _ contracts.HTTPServerAdapter = (*ginHTTPServerAdapter)(nil)

var pError = customerrors.PrefixedError("gin server")

type ginHTTPServerAdapter struct {
	cfg          *config.Config
	engine       *ginVendor.Engine
	logger       loggerContracts.Logger
	routeBuilder contracts.RouteBuilder
	server       *http.Server
}

func New(logger loggerContracts.Logger) contracts.HTTPServerAdapter {
	utils.Assert(logger != nil, pError.New("logger can't be nil"))

	return &ginHTTPServerAdapter{
		logger: logger,
	}
}

func (g *ginHTTPServerAdapter) Setup(cfg *config.Config) error {
	switch cfg.Development {
	case true:
		ginVendor.SetMode(ginVendor.DebugMode)
	default:
		ginVendor.SetMode(ginVendor.ReleaseMode)
	}

	g.cfg = cfg
	g.engine = ginVendor.New()
	g.routeBuilder = NewRouteBuilder(g.engine)

	return nil
}

func (g *ginHTTPServerAdapter) SetupDefaultMiddlewares() {
	g.engine.Use(
		gin.Logger(),
		gin.Recovery(),
		cors.Default(),
		requestid.NewRequestIDMiddleware(),
	)
}

func (s *ginHTTPServerAdapter) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(fmt.Errorf("gin-server: invalid middleware : `%v`", err))
	}

	s.engine.Use(ms...)
}

func (g *ginHTTPServerAdapter) Start(ctx context.Context) error {
	g.logger.Info("starting server...")

	addr := g.cfg.Address()

	g.server = &http.Server{
		Addr:              addr,
		Handler:           g.engine,
		ReadTimeout:       g.cfg.ReadTimeout,
		ReadHeaderTimeout: g.cfg.ReadHeaderTimeout,
		WriteTimeout:      g.cfg.WriteTimeout,
		IdleTimeout:       g.cfg.IdleTimeout,
		MaxHeaderBytes:    g.cfg.MaxHeaderBytes,
	}

	serverErr := make(chan error, 1)
	go func() {
		if err := g.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	g.logger.Infof("server listening on address: %s", addr)

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (g *ginHTTPServerAdapter) Shutdown(ctx context.Context) error {
	g.logger.Info("server shutting down...")

	err := g.server.Shutdown(ctx)

	g.logger.Info("server shutdown done successfully")

	return err
}

func (g *ginHTTPServerAdapter) RouteBuilder() contracts.RouteBuilder {
	return g.routeBuilder
}

func (g *ginHTTPServerAdapter) Type() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.GIN
}
