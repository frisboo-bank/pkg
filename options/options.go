package options

import (
	"errors"
	"fmt"
)

type Option[T any] interface {
	apply(cfg *T) error
}

type OptionFunc[T any] func(*T) error

func (o OptionFunc[T]) apply(cfg *T) error {
	return o(cfg)
}

type OptionBuilder[T any] struct {
	cfg *T
	err error
}

func NewBuilder[T any](base *T) *OptionBuilder[T] {
	return &OptionBuilder[T]{
		cfg: base,
	}
}

func (b *OptionBuilder[T]) With(opts ...Option[T]) *OptionBuilder[T] {
	for _, opt := range opts {
		b.err = errors.Join(b.err, opt.apply(b.cfg))
	}
	return b
}

func (b *OptionBuilder[T]) Build() *T {
	if b.err != nil {
		panic(fmt.Sprintf("invalid configuration: %v", b.err))
	}
	return b.cfg, nil
}

func Apply[T any](base *T) *OptionBuilder[T] {
	return NewBuilder(base)
}
