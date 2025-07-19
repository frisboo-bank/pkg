package httpserver

import (
	"fmt"

	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/gin"
	"frisboo-bank/pkg/http/http_server/options"
	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

func GetInstanceFromOptions(
	options *options.HTTPServerOptions,
	logger loggerContracts.Logger,
) (contracts.HTTPServer, error) {
	instance, err := GetInstance(options.Type, logger)
	if err != nil {
		return nil, err
	}

	return instance.WithOptions(options), nil
}

func GetInstance(
	httpServerType httpservertype.HttpServerType,
	logger loggerContracts.Logger,
) (contracts.HTTPServer, error) {
	switch httpServerType {
	case httpservertype.HttpServerTypes.GIN:
		return gin.NewGinHTTPServer(logger), nil
	default:
		return nil, fmt.Errorf("(http-server-factory) no server of type `%q` exists", httpServerType)
	}
}
