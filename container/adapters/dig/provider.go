package dig

import (
	containerConfig "frisboo-bank/pkg/container/config"
	containerContracts "frisboo-bank/pkg/container/contracts"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/waiter"
	waiterConfig "frisboo-bank/pkg/waiter/config"
)

func ProviderFunc(
	cfg *containerConfig.Config,
	logger loggerContracts.Logger,
) (containerContracts.ContainerAdapter, error) {
	waiterCfg, err := waiterConfig.New(waiterConfig.CancelOnShutdownSignal(true))
	if err != nil {
		return nil, err
	}

	return New(
		cfg,
		waiter.New(waiterCfg, logger),
		logger,
	), nil
}
