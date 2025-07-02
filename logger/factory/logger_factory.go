package factory

import (
	"fmt"
	"frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/logger/logrus"
	"frisboo-bank/pkg/logger/noop"
	"frisboo-bank/pkg/logger/options"
)

func GetInstance(config *options.LogOptions, configs ...options.LogOption) (contracts.Logger, error) {
	for _, c := range configs {
		c(config)
	}

	if config.Type == "" {
		return nil, fmt.Errorf("logger: no logger type specified")
	}

	switch config.Type {
	case options.TypeNoop:
		return noop.NewNoopLogger(), nil
	case options.TypeLogrus:
		return logrus.NewLogrusLogger(config), nil
	}

	return nil, fmt.Errorf("logger: no logger of type `%s` exists", config.Type)
}
