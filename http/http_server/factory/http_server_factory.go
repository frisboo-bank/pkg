package factory

import (
	"fmt"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/gin"
	"frisboo-bank/pkg/http/http_server/options"

	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
)

func GetInstance(config *options.HttpServerOptions, configs ...options.HttpServerOption) (contracts.HttpServer, error) {
	for _, c := range configs {
		c(config)
	}

	switch config.Type {
	case httpservertype.HttpServerTypes.GIN:
		return gin.NewGinHttpServer(config), nil
	default:
		return nil, fmt.Errorf("(http-server-factory) no server of type `%q` exists", config.Type)
	}
}
