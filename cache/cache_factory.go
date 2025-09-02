package cache

import (
	inmemory "frisboo-bank/pkg/cache/adapters/in_memory"
	"frisboo-bank/pkg/cache/adapters/redis"
	"frisboo-bank/pkg/cache/config"
	"frisboo-bank/pkg/cache/contracts"
	"frisboo-bank/pkg/syserrors"

	cachetype "frisboo-bank/pkg/cache/contracts/enums/cache_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
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
	case cachetype.CacheTypes.IN_MEMORY:
		adapter = inmemory.New(logger)
	case cachetype.CacheTypes.REDIS:
		adapter = redis.New(cfg, logger)
	default:
		return nil, NoCacheOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
