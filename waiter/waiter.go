package waiter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"
	"time"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/waiter/contracts"
	waiterOptions "frisboo-bank/pkg/waiter/options"

	"golang.org/x/sync/errgroup"
)

type waiter struct {
	ctx        context.Context
	hooks      []contracts.WaiterHook
	cancel     context.CancelFunc
	isWaiting  bool
	mu         sync.Mutex
	waitOnce   sync.Once
	cancelOnce sync.Once
	logger     loggerContracts.Logger
}

func NewWaiter(options ...waiterOptions.WaiterOption) contracts.Waiter {
	cfg := waiterOptions.GetOptionsWithDefault(options...)

	w := &waiter{}

	parentContext := cfg.ParentContext
	if parentContext == nil {
		parentContext = context.Background()
	}

	w.ctx, w.cancel = context.WithCancel(parentContext)

	if cfg.CancelOnShutdownSignal {
		var signalCancel context.CancelFunc

		signalCtx, signalCancel := signal.NotifyContext(
			w.ctx,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		w.ctx = signalCtx
		parentCancel := w.cancel

		w.cancel = func() {
			signalCancel()
			parentCancel()
		}
	}

	return w
}

func (w *waiter) Wait() error {
	var err error

	w.waitOnce.Do(func() {
		w.mu.Lock()
		w.isWaiting = true
		hooks := slices.Clone(w.hooks)
		w.mu.Unlock()

		defer w.cancel()

		group, gctx := errgroup.WithContext(w.ctx)

		for _, hook := range hooks {
			waitFn := hook.Wait
			cleanupFn := hook.Cleanup

			if waitFn != nil {
				group.Go(func() error {
					waitCtx, waitCancel := context.WithTimeout(gctx, 200*time.Millisecond)
					defer waitCancel()

					return waitFn(waitCtx)
				})
			}

			if cleanupFn != nil {
				group.Go(func() error {
					<-gctx.Done()

					cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
					defer cleanupCancel()

					if err := cleanupFn(cleanupCtx); err != nil {
						fmt.Printf("waiter: failed to cleanup with error: %v", err)
					}

					return nil
				})
			}
		}

		err = group.Wait()
	})

	return err
}

func (w *waiter) Add(hooks ...contracts.WaiterHook) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isWaiting {
		panic(fmt.Errorf("waiter: can't call Add() after Wait() was called"))
	}

	w.hooks = append(w.hooks, hooks...)
}

func (w *waiter) Cancel() {
	w.cancelOnce.Do(func() {
		if w.cancel == nil {
			return
		}

		w.cancel()
	})
}

func (w *waiter) Context() context.Context {
	return w.ctx
}

func (w *waiter) Logger() loggerContracts.Logger {
	return w.logger
}
