package redis

import (
	"context"
	"time"

	"frisboo-bank/pkg/cache/config"
	"frisboo-bank/pkg/cache/contracts"
	"frisboo-bank/pkg/syserrors"

	cachetype "frisboo-bank/pkg/cache/contracts/enums/cache_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"github.com/redis/go-redis/v9"
)

var _ contracts.CacheAdapter = (*redisClientAdapter)(nil)

type redisClientAdapter struct {
	cfg    *config.Config
	ctx    context.Context
	client redis.UniversalClient
	logger loggerContracts.Logger
}

func New(cfg *config.Config, logger loggerContracts.Logger) contracts.CacheAdapter {
	syserrors.AssertNotNil("cfg", cfg)
	syserrors.AssertNotNil("logger", logger)

	return &redisClientAdapter{
		logger: logger,
		cfg:    cfg,
		ctx:    nil,
		client: nil,
	}
}

func (r *redisClientAdapter) Close(ctx context.Context) error {
	panic("unimplemented")
}

func (r *redisClientAdapter) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	panic("unimplemented")
}

func (r *redisClientAdapter) Delete(ctx context.Context, key string) error {
	panic("unimplemented")
}

func (r *redisClientAdapter) Exists(ctx context.Context, key string) (bool, error) {
	panic("unimplemented")
}

func (r *redisClientAdapter) Flush(ctx context.Context) error {
	panic("unimplemented")
}

func (r *redisClientAdapter) Get(ctx context.Context, key string, dest any) (found bool, err error) {
	panic("unimplemented")
}

func (r *redisClientAdapter) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	panic("unimplemented")
}

func (r *redisClientAdapter) Logger() loggerContracts.Logger {
	panic("unimplemented")
}

func (r *redisClientAdapter) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	panic("unimplemented")
}

func (r *redisClientAdapter) Type() cachetype.CacheType {
	return cachetype.CacheTypes.REDIS
}
