package di

import (
	"frisboo-bank/pkg/di"
)

var LoggerFactory = func(c di.Container) (any, error) {
	return nil, nil

	// config, err := c.Get("config")
	// if err != nil {
	// 	return nil, err
	// }

	// loggerConfig, err := config.ProvideLogConfig(config)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// switch loggerConfig.Type {
	// case logger.LogTypeLogrus:
	// 	return logrus.NewLogrusLogger(loggerConfig), nil
	// }
	//
	// return nil, fmt.Errorf("Logger: no logger of type: %s", loggerConfig.Level)
}
