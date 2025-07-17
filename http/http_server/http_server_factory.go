package httpserver

import (
	"fmt"
	"frisboo-bank/pkg/http/http_server/contracts"
	"frisboo-bank/pkg/http/http_server/gin"

	httpservertype "frisboo-bank/pkg/http/http_server/options/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

func GetInstance(
	httpServerType httpservertype.HTTPServerType,
	logger loggerContracts.Logger,
) (contracts.HTTPServer, error) {
	switch httpServerType {
	case httpservertype.HTTPServerTypes.GIN:
		return gin.NewGinHTTPServer(logger), nil
	default:
		return nil, fmt.Errorf("(http-server-factory) no server of type `%q` exists", httpServerType)
	}
}
