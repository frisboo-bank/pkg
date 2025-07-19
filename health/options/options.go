package options

import (
	"net/http"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"
)

// default options
const (
	EndpointPath string = "/healthz"

	StatusTypeUp   contracts.StatusType = "Up"
	StatusTypeDown contracts.StatusType = "Down"

	StatusCodeUp   int = http.StatusOK
	StatusCodeDown int = http.StatusServiceUnavailable
)

type HealthOptions struct {
	EndpointPath   string `mapstructure:"endpointPath"`
	StatusCodeUp   int    `mapstructure:"statusCodeUp"`
	StatusUp       string `mapstructure:"statusUp"`
	StatusCodeDown int    `mapstructure:"statusCodeDown"`
	StatusDown     string `mapstructure:"statusDown"`
}

func ProvideHealthOptions(loader configContracts.ConfigLoader, env environment.Environment) (*HealthOptions, error) {
	return config.LoadOptions[HealthOptions](loader, env)
}
