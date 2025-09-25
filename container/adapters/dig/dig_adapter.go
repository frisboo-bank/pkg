package dig

import (
	"context"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/validation"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

var _ contracts.ContainerAdapter = (*digAdapter)(nil)

type DigAdapterConfig struct{}

type digAdapter struct {
	cfg    *config.Config
	dig    *dig.Container
	logger loggerContracts.Logger
	waiter waiterContracts.Waiter
	hooks  []string
}

func New(
	cfg *config.Config,
	logger loggerContracts.Logger,
	waiter waiterContracts.Waiter,
) (contracts.ContainerAdapter, error) {
	validation.AssertNotNil("cfg", cfg)
	validation.AssertNotNil("waiter", waiter)
	validation.AssertNotNil("logger", logger)

	return &digAdapter{
		cfg:    cfg,
		dig:    dig.New(),
		waiter: waiter,
		logger: logger,
	}, nil
}

func (a *digAdapter) Start(_ context.Context) error {
	hooks, err := a.resolveHooks()
	if err != nil {
		return err
	}
	if err := a.waiter.AddHooks(hooks...); err != nil {
		return err
	}
	return a.waiter.Wait()
}

func (a *digAdapter) Stop(_ context.Context) error {
	a.waiter.Cancel()
	return nil
}

func (a *digAdapter) Type() containertype.ContainerType {
	return containertype.ContainerTypes.DIG
}
