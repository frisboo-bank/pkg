package gin

import (
	"context"
	"errors"
	"net/http"

	requestid "frisboo-bank/pkg/http/http_server/adapters/gin/middlewares/request_id"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	"github.com/gin-contrib/cors"
	ginVendor "github.com/gin-gonic/gin"
)

var _ contracts.HTTPServerAdapter = (*ginHTTPServerAdapter)(nil)

type ginHTTPServerAdapter struct {
	cfg          *config.HTTPServerConfig
	engine       *ginVendor.Engine
	logger       loggerContracts.Logger
	routeBuilder contracts.RouteBuilder
	server       *http.Server
}

func New(cfg *config.HTTPServerConfig, logger loggerContracts.Logger) contracts.HTTPServerAdapter {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	switch cfg.Debug {
	case true:
		ginVendor.SetMode(ginVendor.DebugMode)
	default:
		ginVendor.SetMode(ginVendor.ReleaseMode)
	}

	engine := ginVendor.New()

	return &ginHTTPServerAdapter{
		cfg:          cfg,
		engine:       engine,
		logger:       logger,
		routeBuilder: NewRouteBuilder(engine),
	}
}

func (g *ginHTTPServerAdapter) SetupDefaultMiddlewares() {
	g.engine.Use(
		ginVendor.Logger(),
		ginVendor.Recovery(),
		cors.Default(),
		requestid.NewRequestIDMiddleware(),
	)
}

func (s *ginHTTPServerAdapter) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(syserrors.Wrap(err, "invalid middleware"))
	}

	s.engine.Use(ms...)
}

func (g *ginHTTPServerAdapter) Start(ctx context.Context) error {
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

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (g *ginHTTPServerAdapter) Stop(ctx context.Context) error {
	err := g.server.Shutdown(ctx)
	return err
}

func (g *ginHTTPServerAdapter) ListRoutes() []any {
	panic("unimplemented")
}

func (g *ginHTTPServerAdapter) RouteBuilder() contracts.RouteBuilder {
	return g.routeBuilder
}

func (g *ginHTTPServerAdapter) Logger() loggerContracts.Logger {
	return g.logger
}

func (g *ginHTTPServerAdapter) Type() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.GIN
}
