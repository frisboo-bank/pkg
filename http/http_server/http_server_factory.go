package httpserver

import (
	"fmt"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/gin"

	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

func GetInstanceFromConfig(config *config.HTTPServerConfig, logger loggerContracts.Logger) (contracts.HTTPServer, error) {
	instance, err := GetInstance(config.Type, logger)
	if err != nil {
		return nil, err
	}

	return instance.WithConfig(config), nil
}

func GetInstance(httpServerType httpservertype.HttpServerType, logger loggerContracts.Logger) (contracts.HTTPServer, error) {
	switch httpServerType {
	case httpservertype.HttpServerTypes.GIN:
		return gin.New(logger), nil
	default:
		return nil, fmt.Errorf("(http-server-factory) no server of type `%q` exists", httpServerType)
	}
}
