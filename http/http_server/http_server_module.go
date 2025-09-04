package httpserver

import (
	"context"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/syserrors"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

func ServerFailedToStartError(err error) error {
	return syserrors.Wrap(err, "server failed to start with error", "HTTPServer", "Hooks", "Start")
}

func ServerFailedToStopError(err error) error {
	return syserrors.Wrap(err, "server failed to stop", "HTTPServer", "Hooks", "Stop")
}

func ModuleFunc() module.Module {
	return module.ModuleFunc(
		"http_server",

		provider.ProvideFunc(config.Load),
		provider.ProvideFunc(GetInstance),

		hook.HooksFunc(httpServerHookStart, httpServerHookStop),
	)
}

func httpServerHookStart(httpServer contracts.HTTPServer) waiterContracts.WaitFunc {
	return func(ctx context.Context) error {
		httpServer.SetupDefaultMiddlewares()

		var err error
		go func() {
			err = httpServer.Start(ctx)
		}()

		if err != nil {
			return ServerFailedToStartError(err)
		}

		return nil
	}
}

func httpServerHookStop(httpServer contracts.HTTPServer) waiterContracts.CleanupFunc {
	return func(ctx context.Context) error {
		if err := httpServer.Stop(ctx); err != nil {
			return ServerFailedToStopError(err)
		}

		return nil
	}
}
