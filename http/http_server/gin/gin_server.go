package gin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"frisboo-bank/pkg/http/http_server/contracts"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	requestid "frisboo-bank/pkg/http/http_server/gin/middlewares/request_id"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ginHTTPServer struct {
	contracts.BaseHTTPServer
	engine     *gin.Engine
	httpServer *http.Server
}

func New(logger loggerContracts.Logger) contracts.HTTPServer {
	server := &ginHTTPServer{}
	server.Init(server)

	server.engine = gin.New()
	server.WithRouteBuilder(NewRouteBuilder(server.engine))
	server.WithLogger(logger)

	return server
}

func (s *ginHTTPServer) SetupInstance() {
	switch s.Development() {
	case true:
		gin.SetMode(gin.DebugMode)
	default:
		gin.SetMode(gin.ReleaseMode)
	}

	s.httpServer = &http.Server{
		Addr:              s.Address(),
		Handler:           s.engine,
		ReadTimeout:       s.ReadTimeout(),
		ReadHeaderTimeout: s.ReadHeaderTimeout(),
		WriteTimeout:      s.WriteTimeout(),
		IdleTimeout:       s.IdleTimeout(),
		MaxHeaderBytes:    s.MaxHeaderBytes(),
	}
}

func (s *ginHTTPServer) Start() error {
	s.SetupInstance()

	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *ginHTTPServer) Shutdown(ctx context.Context) error {
	if s.httpServer == nil {
		return fmt.Errorf("gin-server: looks like there is no server running")
	}

	return s.httpServer.Shutdown(ctx)
}

func (s *ginHTTPServer) AddMiddlewares(middlewares ...any) {
	ms, err := ToMiddlewaresType(middlewares...)
	if err != nil {
		panic(fmt.Errorf("gin-server: invalid middleware : `%v`", err))
	}

	s.engine.Use(ms...)
}

func (s *ginHTTPServer) SetupDefaultMiddlewares() {
	s.AddMiddlewares(
		gin.Logger(),
		gin.Recovery(),
		cors.Default(),
		requestid.NewRequestIDMiddleware(),
	)
}

func (s *ginHTTPServer) Instance() any {
	return s.engine
}

func (s *ginHTTPServer) Type() httpservertype.HttpServerType {
	return httpservertype.HttpServerTypes.GIN
}
