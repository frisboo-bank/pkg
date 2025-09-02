package httpserver

import (
	"frisboo-bank/pkg/http/http_server/adapters/gin"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/http/http_server/contracts"
	httpservertype "frisboo-bank/pkg/http/http_server/contracts/enums/http_server_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
)

func NoServerOfTypeError(sType httpservertype.HttpServerType) error {
	return syserrors.Newf("no http server of type `%q` exists", sType)
}

func GetInstance(cfg *config.Config, logger loggerContracts.Logger) (contracts.HTTPServer, error) {
	var adapter contracts.HTTPServerAdapter

	switch cfg.Type {
	case httpservertype.HttpServerTypes.GIN:
		adapter = gin.New(cfg, logger)
	default:
		return nil, NoServerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
