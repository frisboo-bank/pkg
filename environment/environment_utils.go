package environment

import (
	"fmt"
	"strings"

	"frisboo-bank/pkg/constants"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func NewEnvironment(env string) (Environment, error) {
	return ParseEnvironment(env)
}

func (e Environment) IsEnvironment(env Environment) bool {
	return e == env
}

func (e Environment) IsDevelopment() bool {
	return e == Environments.DEVELOPMENT
}

func (e Environment) IsTesting() bool {
	return e == Environments.TESTING
}

func (e Environment) IsPreprod() bool {
	return e == Environments.PREPROD
}

func (e Environment) IsProduction() bool {
	return e == Environments.PRODUCTION
}

func Load(fallback ...Environment) (Environment, error) {
	environment := Environments.DEVELOPMENT
	if len(fallback) > 0 {
		environment = fallback[0]
	}

	viper.AutomaticEnv()

	_ = godotenv.Load()

	manualEnv := strings.TrimSpace(viper.GetString(constants.APP_ENV))

	if manualEnv != "" {
		var err error
		environment, err = NewEnvironment(manualEnv)
		if err != nil {
			return Environment{}, err
		}
	}

	fmt.Printf("Environment: %s\n", environment)

	return environment, nil
}
