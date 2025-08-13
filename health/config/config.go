package config

import (
	"net/http"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/health/contracts"

	configContracts "frisboo-bank/pkg/config/contracts"
)

// default options
const (
	EndpointPath   string               = "/healthz"
	StatusTypeUp   contracts.StatusType = "Up"
	StatusTypeDown contracts.StatusType = "Down"
	StatusCodeUp   int                  = http.StatusOK
	StatusCodeDown int                  = http.StatusServiceUnavailable
)

type HealthConfig struct {
	EndpointPath   string `mapstructure:"endpointPath"`
	StatusCodeUp   int    `mapstructure:"statusCodeUp"`
	StatusUp       string `mapstructure:"statusUp"`
	StatusCodeDown int    `mapstructure:"statusCodeDown"`
	StatusDown     string `mapstructure:"statusDown"`
}

func ProvideHealthConfig(loader configContracts.ConfigLoader, env environment.Environment) (*HealthConfig, error) {
	return config.LoadConfig[HealthConfig](loader, env, "health")
}
