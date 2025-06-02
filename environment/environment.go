package environment

import (
	"frisboo-bank/pkg/constants"
)

type Environment string

var (
	Development = Environment(constants.DEVELOPMENT)
	Testing     = Environment(constants.TESTING)
	PreProd     = Environment(constants.PREPROD)
	Production  = Environment(constants.PRODUCTION)
)

func (e Environment) IsEnvironment(env Environment) bool {
	return e == env
}

func (e Environment) IsDevelopment() bool {
	return e.IsEnvironment(Development)
}

func (e Environment) IsTesting() bool {
	return e.IsEnvironment(Testing)
}

func (e Environment) IsPreprod() bool {
	return e.IsEnvironment(PreProd)
}

func (e Environment) IsProduction() bool {
	return e.IsEnvironment(Production)
}

func (e Environment) ToString() string {
	return string(e)
}
