package redis

import (
	"context"
	"time"

	"frisboo-bank/pkg/cache/config"
	"frisboo-bank/pkg/cache/contracts"
	customErrors "frisboo-bank/pkg/custom_errors"
	"frisboo-bank/pkg/utils"

	cachetype "frisboo-bank/pkg/cache/contracts/enums/cache_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"github.com/redis/go-redis/v9"
)

var _ contracts.CacheAdapter = (*redisClientAdapter)(nil)

var pError = customErrors.PrefixedError("redis server")

type redisClientAdapter struct {
	cfg    *config.Config
	ctx    context.Context
	client redis.UniversalClient
	logger loggerContracts.Logger
}

func New(logger loggerContracts.Logger) contracts.CacheAdapter {
	syserrors.Assert(logger != nil, pError.New("logger can't be nil"))

	return &redisClientAdapter{
		logger: logger,
	}
}

func (r *redisClientAdapter) Setup(cfg *config.Config) error {
	r.cfg = cfg
	r.client = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})

	return nil
}

func (r *redisClientAdapter) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *redisClientAdapter) Flush() error {
	return r.client.FlushAll(r.ctx).Err()
}

func (r *redisClientAdapter) Get(key string) any {
	panic("unimplemented")
}

func (r *redisClientAdapter) Has(key string) bool {
	panic("unimplemented")
}

func (r *redisClientAdapter) Set(key string, value any, t time.Duration) error {
	return r.client.Set(r.ctx, key, value, t).Err()
}

func (r *redisClientAdapter) Type() cachetype.CacheType {
	return cachetype.CacheTypes.REDIS
}
