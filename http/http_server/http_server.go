package httpserver

import (
	"context"

	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/syserrors"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

var _ contracts.HTTPServer = (*httpServer)(nil)

type httpServer struct {
	adapter contracts.HTTPServerAdapter
}

func New(adapter contracts.HTTPServerAdapter) contracts.HTTPServer {
	syserrors.Assert(adapter != nil, "adapter can't be nil")
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

func (h *httpServer) Shutdown(ctx context.Context) error {
	return h.adapter.Shutdown(ctx)
}

func (h *httpServer) Type() httpservertype.HttpServerType {
	return h.adapter.Type()
}

func (h *httpServer) RouteBuilder() contracts.RouteBuilder {
	return h.adapter.RouteBuilder()
}

func (h *httpServer) Logger() loggerContracts.Logger {
	return h.adapter.Logger()
}
