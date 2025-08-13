package logger

import (
	"fmt"

	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
	"frisboo-bank/pkg/logger/logrus"
	"frisboo-bank/pkg/logger/noop"
	"frisboo-bank/pkg/logger/zerolog"
	"frisboo-bank/pkg/options"
)

func GetInstance(lType loggertype.LoggerType, opt *options.OptionBuilder[config.Config]) (contracts.Logger, error) {
	var adapter contracts.LoggerAdapter

	switch lType {
	case loggertype.LoggerTypes.LOGRUS:
		adapter = logrus.New()
	case loggertype.LoggerTypes.NOOP:
		adapter = noop.New()
	case loggertype.LoggerTypes.ZEROLOG:
		adapter = zerolog.New()
	default:
		return nil, fmt.Errorf("logger-factory: type %q not supported", lType)
	}

	log, err := New(adapter, opt)
	if err != nil {
		return nil, err
	}

	return log, nil
}
