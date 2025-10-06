package httpserver

import (
	"frisboo-bank/pkg/http/http_server/adapters/echo"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
)

func NoHTTPServerOfTypeError(name string, sType httpservertype.HttpServerType) error {
	return syserrors.Newf("http-server type %s for server %s does not exist", sType, name)
}

func GetInstance(name string, cfg *config.Config, logger loggerContracts.Logger) (contracts.HTTPServer, error) {
	validation.AssertNotEmpty("name", name)
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	var adapter contracts.HTTPServerAdapter
	switch cfg.Type {
	case httpservertype.HttpServerTypes.ECHO:
		adapter = echo.New(name, cfg, logger, nil)
	default:
		return nil, NoHTTPServerOfTypeError(name, cfg.Type)
	}

	return New(adapter), nil
}
