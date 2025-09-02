package logger

import (
	"frisboo-bank/pkg/logger/adapters/logrus"
	"frisboo-bank/pkg/logger/adapters/noop"
	"frisboo-bank/pkg/logger/adapters/zerolog"
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"

	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

func NoContainerOfTypeError(sType loggertype.LoggerType) error {
	return syserrors.Newf("no container of type %q exists", sType)
}

func GetInstance(cfg *config.Config) (contracts.Logger, error) {
	var adapter contracts.LoggerAdapter
	switch cfg.Type {
	case loggertype.LoggerTypes.LOGRUS:
		adapter = logrus.New(cfg)
	case loggertype.LoggerTypes.NOOP:
		adapter = noop.New(cfg)
	case loggertype.LoggerTypes.ZEROLOG:
		adapter = zerolog.New(cfg)
	default:
		return nil, NoContainerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
