package builder

import (
	"frisboo-bank/pkg/application/contracts"
)

type applicationInfrastructure struct {
	contracts.Application
}

var _ contracts.ApplicationInfrastructure = (*applicationInfrastructure)(nil)

func NewApplicationInfrastructure(app contracts.Application) contracts.ApplicationInfrastructure {
	return &applicationInfrastructure{app}
}

func (a *applicationInfrastructure) ConfigureInfrastructure() error {
	return nil
}
