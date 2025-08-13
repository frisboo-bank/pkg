package health

import (
	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/health/config"
	"frisboo-bank/pkg/health/contracts"
	httpServerContracts "frisboo-bank/pkg/http/http_server/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/utils"

	"github.com/gin-gonic/gin"
)

var _ contracts.HealthEndpoint = (*healthEndpoint)(nil)

var pHEError = customerrors.PrefixedError("health endpoint")

type healthEndpoint struct {
	cfg           *config.Config
	healthService contracts.HealthService
	httpServer    httpServerContracts.HTTPServer
}

func NewHealthEndpoint(
	logger loggerContracts.Logger,
	httpServer httpServerContracts.HTTPServer,
	healthService contracts.HealthService,
	opts *options.OptionBuilder[config.Config],
) contracts.HealthEndpoint {
	utils.Assert(logger != nil, pHEError.New("logger can't be nil"))
	utils.Assert(httpServer != nil, pHEError.New("httpServer can't be nil"))
	utils.Assert(healthService != nil, pHEError.New("healthService can't be nil"))
	utils.Assert(opts != nil, pHEError.New("opts can't be nil"))

	cfg := opts.Build()

	return &healthEndpoint{
		cfg:           cfg,
		healthService: healthService,
		httpServer:    httpServer,
	}
}

func (e *healthEndpoint) RegisterEndpoints() {
	e.httpServer.RouteBuilder().RegisterRoutes(func(server any) {
		server.(*gin.Engine).GET(e.cfg.EndpointPath, e.checkHealth)
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
