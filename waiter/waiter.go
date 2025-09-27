package waiter

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/validation"
	"frisboo-bank/pkg/waiter/config"
	"frisboo-bank/pkg/waiter/contracts"

	"golang.org/x/sync/errgroup"
)

var _ contracts.Waiter = (*waiter)(nil)

type waiter struct {
	cfg   *config.Config
	hooks map[string]contracts.WaiterHook

	cancel     context.CancelFunc
	cancelOnce sync.Once
	ctx        context.Context

	isWaiting bool
	mu        sync.Mutex
	waitOnce  sync.Once

	logger loggerContracts.Logger
}

func New(logger loggerContracts.Logger, opts ...config.Option) (contracts.Waiter, error) {
	validation.AssertNotNil("logger", logger)

	cfg, err := config.New(opts...)
	if err != nil {
		return nil, err
	}

	parentCtx := cfg.ParentContext
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	ctx, cancel := context.WithCancel(parentCtx)

	if cfg.CancelOnShutdownSignal {
		signalCtx, signalCancel := signal.NotifyContext(ctx,
			os.Interrupt,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)

		parentCancel := cancel
		cancel = func() {
			signalCancel()
			parentCancel()
		}
		ctx = signalCtx
	}

	w := &waiter{
		cfg:    &cfg,
		logger: logger,
		cancel: cancel,
		ctx:    ctx,
		hooks:  make(map[string]contracts.WaiterHook),
	}

	return w, nil
}

func (w *waiter) AddHooks(hooks ...contracts.WaiterHook) error {
	for _, h := range hooks {
		if err := w.AddHook(h); err != nil {
			return err
		}
	}
	return nil
}

func (w *waiter) AddHook(hook contracts.WaiterHook) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isWaiting {
		return syserrors.New("waiter: can't call Add() after Wait() was called")
	}

	if hook.Name == "" {
		return syserrors.New("waiter: hook has no name")
	}

	if _, exists := w.hooks[hook.Name]; exists {
		return syserrors.Newf("waiter: a hook %s was already registered", hook.Name)
	}

	w.hooks[hook.Name] = hook
	return nil
}

func (w *waiter) Wait() error {
	var err error
	w.waitOnce.Do(func() {
		err = w.run()
	})
	return err
}

func (w *waiter) run() error {
	w.mu.Lock()
	w.isWaiting = true
	w.mu.Unlock()

	defer w.cancel()

	group, gCtx := errgroup.WithContext(w.ctx)

	for _, hook := range w.hooks {
		waitFn := hook.Wait
		cleanupFn := hook.Cleanup

		if waitFn != nil {
			group.Go(func() error {
				waitCtx := gCtx

				if w.cfg.WaitTimeout > 0 {
					var cancel context.CancelFunc
					waitCtx, cancel = context.WithTimeout(gCtx, w.cfg.WaitTimeout)
					defer cancel()
				}

				return waitFn(waitCtx)
			})
		}

		if cleanupFn != nil {
			group.Go(func() error {
				<-gCtx.Done()

				cleanupCtx, cancel := context.WithTimeout(context.Background(), w.cfg.CleanupTimeout)
				defer cancel()

				if err := cleanupFn(cleanupCtx); err != nil {
					w.logger.Errorf("hook cleanup failed with error: %v", err)
				}
				return nil
			})
		}
	}

	return group.Wait()
}

func (w *waiter) Cancel() {
	w.cancelOnce.Do(func() {
		if w.cancel != nil {
			w.cancel()
		}
	})
}
