package environment

import (
	"fmt"
	"strings"

	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/syserrors"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Environment string

var (
	Development = Environment(constants.DEVELOPMENT)
	Testing     = Environment(constants.TESTING)
	PreProd     = Environment(constants.PREPROD)
	Production  = Environment(constants.PRODUCTION)
)

func NewEnvironment(env string) Environment {
	switch env {
	case constants.DEVELOPMENT:
		return Development
	case constants.TESTING:
		return Testing
	case constants.PREPROD:
		return PreProd
	case constants.PRODUCTION:
		return Production
	default:
		panic(syserrors.Newf("environment: env `%s` is not a valid", env))
	}
}

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

func Load(fallback ...Environment) Environment {
	environment := Development
	if len(fallback) > 0 {
		environment = fallback[0]
	}

	viper.AutomaticEnv()

	_ = godotenv.Load()

	manualEnv := viper.GetString(constants.APP_ENV)
	if trimmed := strings.TrimSpace(manualEnv); trimmed != "" {
		environment = NewEnvironment(trimmed)
	}

	fmt.Printf("Environment: %s\n", environment)

	return environment
}
