package httpserver

import (
	"context"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/options"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

var _ contracts.HTTPServer = (*httpServer)(nil)

type httpServer struct {
	cfg     *config.Config
	adapter contracts.HTTPServerAdapter
	logger  loggerContracts.Logger
}

func New(
	adapter contracts.HTTPServerAdapter,
	logger loggerContracts.Logger,
	opt *options.OptionBuilder[config.Config],
) (contracts.HTTPServer, error) {
	cfg := opt.Build()

	server := &httpServer{
		cfg:     cfg,
		adapter: adapter,
		logger:  logger,
	}

	if err := adapter.Setup(cfg); err != nil {
		return nil, err
	}

	return server, nil
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
	return h.logger
}
