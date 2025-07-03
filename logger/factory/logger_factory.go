package factory

import (
	"fmt"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/logrus"
	"frisboo-bank/pkg/logger/noop"
	"frisboo-bank/pkg/logger/options"

	logtype "frisboo-bank/pkg/logger/options/enums/log_type"
)

func GetInstance(config *options.LogOptions, configs ...options.LogOption) (contracts.Logger, error) {
	for _, c := range configs {
		c(config)
	}

	switch config.Type {
	case logtype.LogTypes.LOGRUS:
		return logrus.NewLogrusLogger(config), nil
	case logtype.LogTypes.NOOP:
		return noop.NewNoopLogger(), nil
	default:
		return nil, fmt.Errorf("logger-factory: type %q not supported", config.Type)
	}
}
