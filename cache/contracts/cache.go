package contracts

import (
	"context"
	"time"

	cachetype "frisboo-bank/pkg/cache/contracts/enums/cache_type"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
)

type (
	DataSerializer interface {
		Marshall(value any) ([]byte, error)
		Unmarshall(data []byte, value any) error
	}

	Cache interface {
		Set(key string, value any, ttl time.Duration) error
		Get(key string, dest any) (found bool, err error)
		Delete(key string) error
		Exists(key string) (exists bool, err error)
		Increment(key string, delta int64) (newValue int64, err error)
		Decrement(key string, delta int64) (newValue int64, err error)
		Flush() error
		Close() error
		WithContext(ctx context.Context) Cache
		Type() cachetype.CacheType
		Logger() loggerContracts.Logger
	}

	CacheAdapter interface {
		Set(ctx context.Context, key string, value any, ttl time.Duration) error
		Get(ctx context.Context, key string, dest any) (found bool, err error)
		Delete(ctx context.Context, key string) error
		Exists(ctx context.Context, key string) (bool, error)
		Increment(ctx context.Context, key string, delta int64) (int64, error)
		Decrement(ctx context.Context, key string, delta int64) (int64, error)
		Flush(ctx context.Context) error
		Type() cachetype.CacheType
		Close(ctx context.Context) error
		Logger() loggerContracts.Logger
	}
)
