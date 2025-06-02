package waiter

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type (
	WaitFunc     func(ctx context.Context) error
	CleanupFunc  func()
	WaiterOption func(config *waiterConfig)
)

type Waiter interface {
	Add(fns ...WaitFunc)
	Cleanup(fns ...CleanupFunc)
	Wait() error
	Context() context.Context
	CancelFunc() context.CancelFunc
}

type waiter struct {
	ctx          context.Context
	waitFuncs    []WaitFunc
	cleanupFuncs []CleanupFunc
	cancel       context.CancelFunc
}

type waiterConfig struct {
	parentContext context.Context
	catchSignals  bool
}

func ParentContext(ctx context.Context) WaiterOption {
	return func(config *waiterConfig) {
		config.parentContext = ctx
	}
}

func CatchSignals() WaiterOption {
	return func(config *waiterConfig) {
		config.catchSignals = true
	}
}

func NewWaiter(options ...WaiterOption) Waiter {
	cfg := &waiterConfig{
		parentContext: context.Background(),
		catchSignals:  false,
	}

	for _, option := range options {
		option(cfg)
	}

	w := &waiter{
		waitFuncs:    []WaitFunc{},
		cleanupFuncs: []CleanupFunc{},
	}

	w.ctx, w.cancel = context.WithCancel(cfg.parentContext)

	if cfg.catchSignals {
		w.ctx, w.cancel = signal.NotifyContext(w.ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	}

	return w
}

// Add implements Waiter.
func (w *waiter) Add(fns ...WaitFunc) {
	w.waitFuncs = append(w.waitFuncs, fns...)
}

// Cleanup implements Waiter.
func (w *waiter) Cleanup(fns ...CleanupFunc) {
	w.cleanupFuncs = append(w.cleanupFuncs, fns...)
}

// Wait implements Waiter.
func (w *waiter) Wait() error {
	group, ctx := errgroup.WithContext(w.ctx)

	group.Go(func() error {
		<-ctx.Done()
		w.cancel()
		return nil
	})

	for _, fn := range w.waitFuncs {
		group.Go(func() error {
			return fn(ctx)
		})
	}

	for _, fn := range w.cleanupFuncs {
		defer fn()
	}

	return group.Wait()
}

// Context implements Waiter.
func (w *waiter) Context() context.Context {
	return w.ctx
}

// CancelFunc implements Waiter.
func (w *waiter) CancelFunc() context.CancelFunc {
	return w.cancel
}
