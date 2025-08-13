package container

import (
	"fmt"

	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/contracts/enums/container_type"
	"frisboo-bank/pkg/container/dig"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/waiter"
	waiterConfig "frisboo-bank/pkg/waiter/config"
)

func GetInstance(
	cType containertype.ContainerType,
	logger loggerContracts.Logger,
	opt *options.OptionBuilder[config.Config],
) (contracts.Container, error) {
	var adapter contracts.ContainerAdapter

	switch cType {
	case containertype.ContainerTypes.DIG:
		waiterConfig := waiterConfig.Apply().
			With(waiterConfig.CancelOnShutdownSignal(true))

		wait, err := waiter.New(logger, waiterConfig)
		if err != nil {
			return nil, err
		}
		adapter = dig.New(logger, wait)
	default:
		return nil, fmt.Errorf("container-factory: type %q not supported", cType)
	}

	di, err := New(adapter, logger, opt)
	if err != nil {
		return nil, err
	}

	return di, nil
}
