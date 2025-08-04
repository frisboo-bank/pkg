package logger

import (
	"fmt"

	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/logrus"
	"frisboo-bank/pkg/logger/noop"
	"frisboo-bank/pkg/logger/zerolog"

	loggertype "frisboo-bank/pkg/logger/contracts/enums/logger_type"
)

func GetInstanceFromConfig(config *config.LoggerConfig) (contracts.Logger, error) {
	instance, err := GetInstance(config.Type)
	if err != nil {
		return nil, err
	}

	return instance.WithConfig(config), nil
}

func GetInstance(loggerType loggertype.LoggerType) (contracts.Logger, error) {
	switch loggerType {
	case loggertype.LoggerTypes.LOGRUS:
		return logrus.New(), nil
	case loggertype.LoggerTypes.NOOP:
		return noop.New(), nil
	case loggertype.LoggerTypes.ZEROLOG:
		return zerolog.New(), nil
	}

	return nil, fmt.Errorf("logger-factory: type %q not supported", loggerType)
}
