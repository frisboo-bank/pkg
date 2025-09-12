package cache

import (
	"frisboo-bank/pkg/cache/adapters/redis"
	"frisboo-bank/pkg/cache/config"
	"frisboo-bank/pkg/cache/contracts"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
)

func NoCacheOfTypeError(sType cachetype.CacheType) error {
	return syserrors.Newf("no cache of type %q exists", sType)
}

func GetInstance(
	cfg *config.Config,
	codec contracts.DataSerializer,
	logger loggerContracts.Logger,
) (contracts.Cache, error) {
	var adapter contracts.CacheAdapter

	switch cfg.Type {
	case cachetype.CacheTypes.REDIS:
		adapter = redis.New(cfg, logger)
	default:
		return nil, NoCacheOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
