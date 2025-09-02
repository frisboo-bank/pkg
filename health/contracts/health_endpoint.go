package contracts

import loggerContracts "frisboo-bank/pkg/logger/contracts"

type HealthEndpoint interface {
	RegisterEndpoints()
	Logger() loggerContracts.Logger
}
