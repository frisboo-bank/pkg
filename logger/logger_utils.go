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

func GetInstance(cfg *config.Config) (contracts.Logger, error) {
	var adapter contracts.LoggerAdapter

	switch cfg.Type {
	case loggertype.LoggerTypes.LOGRUS:
		adapter = logrus.New(cfg)
	case loggertype.LoggerTypes.ZEROLOG:
		adapter = zerolog.New(cfg)
	default:
		return nil, NoContainerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}

func GetByName(cfgRegistry config.Registry, name string) (contracts.Logger, error) {
	if name == "" {
		return nil, syserrors.New("no logger name specified")
	}
	cfg, err := cfgRegistry.GetByName(name)
	if err != nil {
		return nil, syserrors.Wrapf(err, "failed to load %s config", name)
	}
	log, err := GetInstance(&cfg)
	if err != nil {
		return nil, syserrors.Wrapf(err, "failed to initialize %s logger", name)
	}
	return log, nil
}

func GetByNameWithFallback(cfgRegistry config.Registry, name string, fallback contracts.Logger) (contracts.Logger, error) {
	if name == "" {
		return fallback, nil
	}
	return GetByName(cfgRegistry, name)
}
