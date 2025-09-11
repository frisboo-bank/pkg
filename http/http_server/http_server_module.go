package httpserver

import (
	"fmt"
	"frisboo-bank/pkg/container"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/http/http_server/adapters/gin"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"

	"github.com/davecgh/go-spew/spew"
)

const HTTPServersGroup = "http_servers"

func ServerFailedToStartError(err error) error {
	return syserrors.Wrap(err, "server failed to start with error", "HTTPServer", "Hooks", "Start")
}

func ServerFailedToStopError(err error) error {
	return syserrors.Wrap(err, "server failed to stop", "HTTPServer", "Hooks", "Stop")
}

func ModuleFunc(cfg *config.Config) module.Module {
	validation.AssertNotNil("cfg", cfg)

	m := module.ModuleFunc("http_server",
		provider.ProvideFunc(func() *config.Config { return cfg }),
	)

	for name, sc := range cfg.Servers {
		serverName := fmt.Sprintf("http_server_%s", name)
		serverCfg := sc

		constructor, err := getAdapterConstructor(serverName, serverCfg)
		if err != nil {
			panic(err)
		}

		m.AddProviders(
			provider.ProvideFunc(constructor,
				provider.Name(serverName),
				provider.Group(HTTPServersGroup),
			),
		)
	}

	m.AddHooks(
		hook.HooksFunc(
			container.ConstructorWithParams(httpServerHookStartAll),
			container.ConstructorWithParams(httpServerHookStopAll),
		),
	)

	return m
}

func getAdapterConstructor(name string, cfg *config.Config) (any, error) {
	if err := validation.NotEmpty("Name", name); err != nil {
		return nil, err
	}
	if err := validation.NotNil("cfg", cfg); err != nil {
		return nil, err
	}

	switch cfg.Type {
	case httpservertype.HttpServerTypes.GIN:
		return func(logger loggerContracts.Logger) contracts.HTTPServer {
			adapter := gin.New(cfg, logger)
			return New(adapter)
		}, nil
	}

	return nil, syserrors.Newf("http_server %q is using and invalid type: got %q", name, cfg.Type)
}

func httpServerHookStartAll(params struct {
	Servers []contracts.HTTPServer `group:"http_server"`
},
) {
	spew.Dump(params.Servers)
}

func httpServerHookStopAll(params struct {
	Servers []contracts.HTTPServer `group:"http_server"`
},
) {
	spew.Dump(params.Servers)
}

// func httpServerHookStart(httpServer contracts.HTTPServer) waiterContracts.WaitFunc {
// 	return func(ctx context.Context) error {
// 		httpServer.SetupDefaultMiddlewares()
//
// 		var err error
// 		go func() {
// 			err = httpServer.Start(ctx)
// 		}()
//
// 		if err != nil {
// 			return ServerFailedToStartError(err)
// 		}
//
// 		return nil
// 	}
// }
//
// func httpServerHookStop(httpServer contracts.HTTPServer) waiterContracts.CleanupFunc {
// 	return func(ctx context.Context) error {
// 		if err := httpServer.Stop(ctx); err != nil {
// 			return ServerFailedToStopError(err)
// 		}
//
// 		return nil
// 	}
// }
