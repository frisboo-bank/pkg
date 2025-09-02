package health

import (
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"

	"github.com/gin-gonic/gin"
)

var _ contracts.HealthEndpoint = (*healthEndpoint)(nil)

type healthEndpoint struct {
	cfg           *config.Config
	logger        loggerContracts.Logger
	httpServer    httpServerContracts.HTTPServer
	healthService contracts.HealthService
}

func NewHealthEndpoint(
	cfg *config.Config,
	logger loggerContracts.Logger,
	httpServer httpServerContracts.HTTPServer,
	healthService contracts.HealthService,
) contracts.HealthEndpoint {
	syserrors.AssertNotNil("cfg", cfg)
	syserrors.AssertNotNil("logger", logger)
	syserrors.AssertNotNil("httpServer", httpServer)
	syserrors.AssertNotNil("healthService", healthService)

	return &healthEndpoint{
		cfg:           cfg,
		healthService: healthService,
		httpServer:    httpServer,
		logger:        logger,
	}
}

func (e *healthEndpoint) RegisterEndpoints() {
	e.httpServer.RouteBuilder().RegisterRoutes(func(server any) {
		server.(*gin.Engine).GET(e.cfg.LivenessPath, e.checkHealth)
	})
}

func (e *healthEndpoint) checkHealth(ctx *gin.Context) {
	status := e.healthService.CheckHealth(ctx.Request.Context())
	if !status.IsAllUP() {
		ctx.JSON(e.cfg.StatusCodeDown, status)
		return
	}

	ctx.JSON(e.cfg.StatusCodeUp, status)
}

func (e *healthEndpoint) Logger() loggerContracts.Logger {
	return e.logger
}
