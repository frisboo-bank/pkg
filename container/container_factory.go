package container

import (
	"frisboo-bank/pkg/container/adapters/dig"
	"frisboo-bank/pkg/container/config"
	"frisboo-bank/pkg/container/contracts"
	containertype "frisboo-bank/pkg/container/enums/container_type"
	loggerContracts "frisboo-bank/pkg/logger/contracts"
	"frisboo-bank/pkg/syserrors"
)

func NoContainerOfTypeError(sType containertype.ContainerType) error {
	return syserrors.Newf("no container of type %q exists", sType)
}

func GetInstance(
	cfg *config.Config,
	logger loggerContracts.Logger,
) (contracts.Container, error) {
	var adapter contracts.ContainerAdapter
	var err error

	switch cfg.Type {
	case containertype.ContainerTypes.DIG:
		adapter, err = dig.ProviderFunc(cfg, logger)
	default:
		return nil, NoContainerOfTypeError(cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	return New(adapter), nil
}
