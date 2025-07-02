package endpoints

import (
	"frisboo-bank/pkg/health/contracts"
	"frisboo-bank/pkg/health/options"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"

	"github.com/gin-gonic/gin"
)

type healthCheckEndpoint struct {
	config        *options.HealthOptions
	healthService contracts.HealthService
	httpServer    httpServerContracts.HttpServer
}

var _ contracts.HealthEndpoint = (*healthCheckEndpoint)(nil)

func NewHealthCheckEndpoint(
	config *options.HealthOptions,
	healthService contracts.HealthService,
	httpServer httpServerContracts.HttpServer,
) contracts.HealthEndpoint {
	return &healthCheckEndpoint{
		config:        config,
		healthService: healthService,
		httpServer:    httpServer,
	}
}

func (e *healthCheckEndpoint) RegisterEndpoints() {
	e.httpServer.RouteBuilder().RegisterRoutes(func(server any) {
		server.(*gin.Engine).GET(e.config.EndpointPath, e.checkHealth)
	})
}

func (e *healthCheckEndpoint) checkHealth(ctx *gin.Context) {
	status := e.healthService.CheckHealth(ctx.Request.Context())
	if !status.IsAllUP() {
		ctx.JSON(int(e.config.StatusCodeDown), status)
		return
	}

	ctx.JSON(int(e.config.StatusCodeUp), status)
}
