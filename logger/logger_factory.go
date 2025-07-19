package logger

import (
	"fmt"

	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/logrus"
	"frisboo-bank/pkg/logger/noop"
	"frisboo-bank/pkg/logger/options"
	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

func GetInstanceFromOptions(options *options.LogOptions) (contracts.Logger, error) {
	instance, err := GetInstance(options.Type)
	if err != nil {
		return nil, err
	}

	return instance.WithOptions(options), nil
}

func GetInstance(logType logtype.LogType) (contracts.Logger, error) {
	switch logType {
	case logtype.LogTypes.LOGRUS:
		return logrus.NewLogrusLogger(), nil
	case logtype.LogTypes.NOOP:
		return noop.NewNoopLogger(), nil
	default:
		return nil, fmt.Errorf("logger-factory: type %q not supported", logType)
	}
}
