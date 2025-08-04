package health

import (
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"

	"github.com/gin-gonic/gin"
)

type healthEndpoint struct {
	endpointPath   string
	healthService  contracts.HealthService
	httpServer     httpServerContracts.HTTPServer
	statusCodeDown int
	statusCodeUp   int
}

func (e *healthEndpoint) WithEndpointPath(endpointPath string) contracts.HealthEndpoint {
	e.endpointPath = endpointPath
	return e
}

func (e *healthEndpoint) WithStatusCodeDown(statusCodeDown int) contracts.HealthEndpoint {
	e.statusCodeDown = statusCodeDown
	return e
}

func (e *healthEndpoint) WithStatusCodeUp(statusCodeUp int) contracts.HealthEndpoint {
	e.statusCodeUp = statusCodeUp
	return e
}

var _ contracts.HealthEndpoint = (*healthEndpoint)(nil)

func NewHealthEndpoint(
	healthService contracts.HealthService,
	httpServer httpServerContracts.HTTPServer,
) contracts.HealthEndpoint {
	return &healthEndpoint{
		endpointPath:   config.EndpointPath,
		healthService:  healthService,
		httpServer:     httpServer,
		statusCodeDown: config.StatusCodeDown,
		statusCodeUp:   config.StatusCodeUp,
	}
}

func (e *healthEndpoint) RegisterEndpoints() {
	e.httpServer.RouteBuilder().RegisterRoutes(func(server any) {
		server.(*gin.Engine).GET(e.endpointPath, e.checkHealth)
	})
}

func (e *healthEndpoint) checkHealth(ctx *gin.Context) {
	status := e.healthService.CheckHealth(ctx.Request.Context())
	if !status.IsAllUP() {
		ctx.JSON(e.statusCodeDown, status)
		return
	}

	ctx.JSON(e.statusCodeUp, status)
}
