package logger

import (
	"frisboo-bank/pkg/logger/adapters/logrus"
	"frisboo-bank/pkg/logger/adapters/zerolog"
	"frisboo-bank/pkg/logger/config"
	"frisboo-bank/pkg/logger/contracts"
	loggertype "frisboo-bank/pkg/logger/enums/logger_type"
	"frisboo-bank/pkg/syserrors"
)

func NoContainerOfTypeError(sType loggertype.LoggerType) error {
	return syserrors.Newf("no logger of type %q exists", sType)
}

func GetInstance(name string, cfg config.Config) (contracts.Logger, error) {
	var adapter contracts.LoggerAdapter

	switch cfg.Type {
	case loggertype.LoggerTypes.LOGRUS:
		adapter = logrus.New(&cfg)
	case loggertype.LoggerTypes.ZEROLOG:
		adapter = zerolog.New(&cfg)
	default:
		return nil, NoContainerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}

func GetByName(cfgRegistry config.Registry, name string) (contracts.Logger, error) {
	cfg, err := config.GetConfigByName(cfgRegistry, name)
	if err != nil {
		return nil, err
	}
	return GetInstance(name, cfg)
}

func GetDefault(cfgRegistry config.Registry) (contracts.Logger, error) {
	cfg, err := cfgRegistry.GetDefault()
	if err != nil {
		return nil, err
	}
	return GetInstance("default", cfg)
}

func GetByNameWithFallback(cfgRegistry config.Registry, name string, fallback contracts.Logger) (contracts.Logger, error) {
	if name == "" {
		return fallback, nil
	}
	return GetByName(cfgRegistry, name)
}

func GetByNameOrDefault(cfgRegistry config.Registry, name string) (contracts.Logger, error) {
	if name == "" {
		return GetDefault(cfgRegistry)
	}
	return GetByName(cfgRegistry, name)
}
