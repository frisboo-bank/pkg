package inmemory

import (
	"context"
	"sync"
	"time"

	"frisboo-bank/pkg/cache/config"
	"frisboo-bank/pkg/cache/contracts"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

var _ contracts.CacheAdapter = (*inMemoryAdapter)(nil)

type inMemoryAdapter struct {
	cfg    *config.Config
	logger loggerContracts.Logger
	mu     sync.RWMutex
}

func New(logger loggerContracts.Logger) contracts.CacheAdapter {
	return &inMemoryAdapter{
		logger: logger,
	}
}

func (i *inMemoryAdapter) Close(ctx context.Context) error {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Delete(ctx context.Context, key string) error {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Exists(ctx context.Context, key string) (bool, error) {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Flush(ctx context.Context) error {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Get(ctx context.Context, key string, dest any) (found bool, err error) {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	panic("unimplemented")
}

func (i *inMemoryAdapter) Logger() loggerContracts.Logger {
	return i.logger
}

func (i *inMemoryAdapter) Type() cachetype.CacheType {
	return cachetype.CacheTypes.IN_MEMORY
}
