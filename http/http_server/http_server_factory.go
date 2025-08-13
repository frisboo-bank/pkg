package httpserver

import (
	"fmt"

	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	"frisboo-bank/pkg/http/http_server/gin"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
)

func GetInstance(
	sType httpservertype.HttpServerType,
	logger loggerContracts.Logger,
	opt *options.OptionBuilder[config.Config],
) (contracts.HTTPServer, error) {
	var adapter contracts.HTTPServerAdapter

	switch sType {
	case httpservertype.HttpServerTypes.GIN:
		adapter = gin.New(logger)
	default:
		return nil, fmt.Errorf("(http-server-factory) no server of type `%q` exists", sType)
	}

	server, err := New(adapter, logger, opt)
	if err != nil {
		return nil, err
	}

	return server, nil
}
