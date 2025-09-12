package httpserver

import (
	"frisboo-bank/pkg/container/dependencies/module"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
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

	m := module.ModuleFunc("http_server")

	// for name, sc := range cfg.Servers {
	// 	serverName := fmt.Sprintf("http_server_%s", name)
	// 	serverCfg := sc
	//
	// 	constructor, err := getAdapterConstructor(serverName, serverCfg)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	m.AddProviders(
	// 		provider.ProvideFunc(constructor,
	// 			provider.Name(serverName),
	// 			provider.Group(HTTPServersGroup),
	// 		),
	// 	)
	// }
	//
	// m.AddHooks(
	// 	hook.HooksFunc(
	// 		container.ConstructorWithParams(httpServerHookStartAll),
	// 		container.ConstructorWithParams(httpServerHookStopAll),
	// 	),
	// )

	return m
}

func getAdapterConstructor(name string, cfg *config.Config) (any, error) {
	// if err := validation.NotEmpty("Name", name); err != nil {
	// 	return nil, err
	// }
	// if err := validation.NotNil("cfg", cfg); err != nil {
	// 	return nil, err
	// }
	//
	// switch cfg.Type {
	// case httpservertype.HttpServerTypes.GIN:
	// 	return func(logger loggerContracts.Logger) contracts.HTTPServer {
	// 		adapter := gin.New(cfg, logger)
	// 		return New(adapter)
	// 	}, nil
	// }

	return nil, syserrors.Newf("http_server %q is using and invalid type: got %q", name, cfg.Type)
}

func httpServerHookStartAll(params struct {
	Servers []contracts.HTTPServer `group:"http_server"`
},
) {
}

func httpServerHookStopAll(params struct {
	Servers []contracts.HTTPServer `group:"http_server"`
},
) {
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
