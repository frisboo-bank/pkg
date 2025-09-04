package logger

import (
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	"frisboo-bank/pkg/syserrors"
)

func NoContainerOfTypeError(sType loggertype.LoggerType) error {
	return syserrors.Newf("no container of type %q exists", sType)
}

func GetInstance(cfg *config.Config) (contracts.Logger, error) {
	var adapter contracts.LoggerAdapter
	// switch cfg.Type {
	// case loggertype.LoggerTypes.LOGRUS:
	// 	adapter = logrus.New(cfg)
	// case loggertype.LoggerTypes.ZEROLOG:
	// 	adapter = zerolog.New(cfg)
	// default:
	// 	return nil, NoContainerOfTypeError(cfg.Type)
	// }

	return New(adapter), nil
}
