package database_client

import (
	applicationContracts "frisboo-bank/pkg/application/contracts"
	"frisboo-bank/pkg/database/database_client/config"
	"frisboo-bank/pkg/environment"

	"go.uber.org/dig"

	waiterContracts "frisboo-bank/pkg/waiter/contracts"
)

type DatabaseClientModule struct{}

var _ applicationContracts.Module = (*DatabaseClientModule)(nil)

func (d *DatabaseClientModule) Register(container *dig.Container, waiter waiterContracts.Waiter) error {
	err := container.Provide(func(env environment.Environment) (*config.HttpServerOptions, error) {
		return config.ProvideHttpServerConfig(env)
	})
	if err != nil {
		return err
	}

	return nil
}
