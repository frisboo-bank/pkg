package factory

import (
	"fmt"
	"frisboo-bank/pkg/logger"
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/logrus"
)

func NewInstance(cfg *config.LogOptions) logger.Logger {
	switch cfg.Type {
	case logger.LogTypeLogrus:
		return logrus.NewLogrusLogger(cfg)
	}

	panic(fmt.Errorf("logger: no logger of type: %s exists", cfg.Type))
}
