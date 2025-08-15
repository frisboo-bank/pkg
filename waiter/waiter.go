package waiter

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"frisboo-bank/pkg/customerrors"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/utils"
	"frisboo-bank/pkg/waiter/config"
	"frisboo-bank/pkg/waiter/contracts"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

var pError = customerrors.PrefixedError("waiter")

type waiter struct {
	cancel     context.CancelFunc
	cancelOnce sync.Once
	cfg        *config.Config
	ctx        context.Context
	hooks      []contracts.WaiterHook
	isWaiting  bool
	logger     loggerContracts.Logger
	mu         sync.Mutex
	waitOnce   sync.Once
}

func New(logger loggerContracts.Logger, opts *options.OptionBuilder[config.Config]) (contracts.Waiter, error) {
	utils.Assert(logger != nil, pError.New("logger can't be nil"))
	utils.Assert(opts != nil, pError.New("opts can't be nil"))

	cfg := opts.Build()

	ctx, cancel := context.WithCancel(cfg.ParentContext)

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
	}, nil
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
