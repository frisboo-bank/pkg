package httpserver

import (
	"context"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"
)

var _ contracts.HTTPServer = (*httpServer)(nil)

type httpServer struct {
	adapter contracts.HTTPServerAdapter
}

func New(adapter contracts.HTTPServerAdapter) contracts.HTTPServer {
	validation.AssertNotNil("adapter", adapter)

	return &httpServer{adapter}
}

func (h *httpServer) AddMiddlewares(middlewares ...any) {
	h.adapter.AddMiddlewares(middlewares...)
}

func (h *httpServer) SetupDefaultMiddlewares() {
	h.adapter.SetupDefaultMiddlewares()
}

func (h *httpServer) Start(ctx context.Context) error {
	return h.adapter.Start(ctx)
}

func (h *httpServer) Stop(ctx context.Context) error {
	return h.adapter.Stop(ctx)
}

func (h *httpServer) ListRoutes() []any {
	return h.adapter.ListRoutes()
}

func (h *httpServer) Type() httpservertype.HttpServerType {
	return h.adapter.Type()
}

func (h *httpServer) Config() *config.Config {
	return h.adapter.Config()
}

func (h *httpServer) RouteBuilder() contracts.RouteBuilder {
	return h.adapter.RouteBuilder()
}

func (h *httpServer) Logger() loggerContracts.Logger {
	return h.adapter.Logger()
}
