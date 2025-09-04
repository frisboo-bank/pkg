package cache

import (
	"context"
	"time"

	"frisboo-bank/pkg/cache/contracts"
	cachetype "frisboo-bank/pkg/cache/enums/cache_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"
)

var _ contracts.Cache = (*cache)(nil)

type cache struct {
	ctx     context.Context
	adapter contracts.CacheAdapter
}

func New(adapter contracts.CacheAdapter) contracts.Cache {
	validation.Assert(adapter != nil, "adapter can't be nil")

	return cache{
		adapter: adapter,
	}
}

func (c cache) Close() error {
	panic("unimplemented")
}

func (c cache) Decrement(key string, delta int64) (newValue int64, err error) {
	panic("unimplemented")
}

func (c cache) Delete(key string) error {
	panic("unimplemented")
}

func (c cache) Exists(key string) (exists bool, err error) {
	panic("unimplemented")
}

func (c cache) Flush() error {
	panic("unimplemented")
}

func (c cache) Get(key string, dest any) (found bool, err error) {
	panic("unimplemented")
}

func (c cache) Increment(key string, delta int64) (newValue int64, err error) {
	panic("unimplemented")
}

func (c cache) Logger() loggerContracts.Logger {
	panic("unimplemented")
}

func (c cache) Set(key string, value any, ttl time.Duration) error {
	return c.adapter.Set(c.ctx, key, value, ttl)
}

func (c cache) Type() cachetype.CacheType {
	return c.adapter.Type()
}

func (c cache) WithContext(ctx context.Context) contracts.Cache {
	c.ctx = ctx
	return c
}
