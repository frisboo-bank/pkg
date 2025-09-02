package cache

import (
	"frisboo-bank/pkg/cache/contracts"
	"frisboo-bank/pkg/cache/redis"
	"frisboo-bank/pkg/http/http_server/config"
	"frisboo-bank/pkg/options"

	inmemory "frisboo-bank/pkg/cache/in_memory"

	cachetype "frisboo-bank/pkg/cache/contracts/enums/cache_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

var cError = contracts.ModuleError.WithPrefix("Factory")

func GetInstance(
	cType cachetype.CacheType,
	codec contracts.DataSerializer,
	logger loggerContracts.Logger,
	opt *options.OptionBuilder[config.Config],
) (contracts.Cache, error) {
	var adapter contracts.CacheAdapter

	switch cType {
	case cachetype.CacheTypes.IN_MEMORY:
		adapter = inmemory.New(logger)
	case cachetype.CacheTypes.REDIS:
		adapter = redis.New(logger)
	default:
		return nil, cError.Errorf("no cache of type `%q` exists", cType)
	}

	cache, err := New(adapter, logger, opt)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
