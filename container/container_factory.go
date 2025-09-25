package container

import (
	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
)

func NoContainerOfTypeError(sType containertype.ContainerType) error {
	return syserrors.Newf("no container of type %q exists", sType)
}

func GetInstance(cfg *config.Config, logger loggerContracts.Logger) (contracts.Container, error) {
	var adapter contracts.ContainerAdapter

	switch cfg.Type {
	case containertype.ContainerTypes.DIG:
		// waiterCfg, err := waiterConfig.New(
		// 	waiterConfig.CancelOnShutdownSignal(true),
		// )
		// if err != nil {
		// 	return nil, err
		// }
		//
		// wa := waiter.New(&waiterCfg, logger)
		// adapter, err = dig.New(cfg, logger, wa)
	default:
		return nil, NoContainerOfTypeError(cfg.Type)
	}

	return New(adapter), nil
}
