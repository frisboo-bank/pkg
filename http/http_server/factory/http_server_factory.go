package factory

import (
	"fmt"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/gin"
	"frisboo-bank/pkg/http/http_server/options"
)

func GetInstance(config *options.HttpServerOptions, configs ...options.HttpServerOption) (contracts.HttpServer, error) {
	for _, c := range configs {
		c(config)
	}

	if config.Type == "" {
		return nil, fmt.Errorf("http-server: no server type specified")
	}

	switch config.Type {
	case options.TypeGin:
		return gin.NewGinHttpServer(config), nil
	}

	return nil, fmt.Errorf("http-server: no server of type `%s` exists", config.Type)
}
