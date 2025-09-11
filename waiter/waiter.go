package waiter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	"frisboo-bank/pkg/waiter/config"
	"frisboo-bank/pkg/waiter/contracts"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	cfg config.Config

	cancel     context.CancelFunc
	cancelOnce sync.Once
	ctx        context.Context
	hooks      []contracts.WaiterHook
	isWaiting  bool
	mu         sync.Mutex
	waitOnce   sync.Once

	logger loggerContracts.Logger
}

func New(
	cfg config.Config,
	logger loggerContracts.Logger,
) contracts.Waiter {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("logger", logger)

	parentCtx := cfg.ParentContext
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	if cfg.CancelOnShutdownSignal {
		signalCtx, signalCancel := signal.NotifyContext(
			ctx,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		ctx = signalCtx
		parentCancel := cancel

		cancel = func() {
			signalCancel()
			parentCancel()
		}
	}

	return &waiter{
		cancel: cancel,
		cfg:    cfg,
		ctx:    ctx,
		logger: logger,
	}
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
					waitCtx, waitCancel := context.WithTimeout(gctx, w.cfg.WaitTimeout)
					defer waitCancel()

					return waitFn(waitCtx)
				})
			}

			if cleanupFn != nil {
				group.Go(func() error {
					<-gctx.Done()

					cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), w.cfg.CleanupTimeout)
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
		panic(syserrors.Newf("waiter: can't call Add() after Wait() was called"))
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
