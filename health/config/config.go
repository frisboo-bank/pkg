package config

import (
	"net/http"
	"strings"

	"frisboo-bank/pkg/config"
	configContracts "frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/customerrors"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
)

var pError = customerrors.PrefixedError("health config")

type EnvConfig struct {
	EndpointPath   string `mapstructure:"endpointPath"`
	StatusCodeUp   int    `mapstructure:"statusCodeUp"`
	StatusUp       string `mapstructure:"statusUp"`
	StatusCodeDown int    `mapstructure:"statusCodeDown"`
	StatusDown     string `mapstructure:"statusDown"`
}

func LoadEnvConfig(loader configContracts.ConfigLoader, env environment.Environment) (*EnvConfig, error) {
	return config.LoadConfig[EnvConfig](loader, env, "health")
}

type Config struct {
	EndpointPath   string
	StatusCodeUp   int
	StatusUp       string
	StatusCodeDown int
	StatusDown     string
}

var defaultConfig = &Config{
	EndpointPath:   "/healthz",
	StatusCodeUp:   http.StatusOK,
	StatusUp:       "Up",
	StatusCodeDown: http.StatusServiceUnavailable,
	StatusDown:     "Down",
}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}

func FromEnvConfig(cfg *EnvConfig) *options.OptionBuilder[Config] {
	opts := Apply()

	if cfg.EndpointPath != "" {
		opts.With(EndpointPath(cfg.EndpointPath))
	}

	if cfg.StatusCodeUp != 0 {
		opts.With(StatusCodeUp(cfg.StatusCodeUp))
	}

	if cfg.StatusUp != "" {
		opts.With(StatusUp(cfg.StatusUp))
	}

	if cfg.StatusCodeDown != 0 {
		opts.With(StatusCodeDown(cfg.StatusCodeDown))
	}

	if cfg.StatusDown != "" {
		opts.With(StatusDown(cfg.StatusDown))
	}

	return opts
}

func EndpointPath(endpointPath string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		endpointPath = strings.TrimSpace(endpointPath)

		if endpointPath == "" {
			return pError.New("endpointPath can't be empty")
		}

		cfg.EndpointPath = endpointPath
		return nil
	})
}

func StatusCodeUp(statusCodeUp int) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if statusCodeUp <= 0 {
			return pError.New("statusCodeUp must be positive")
		}

		cfg.StatusCodeUp = statusCodeUp
		return nil
	})
}

func StatusUp(statusUp string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		statusUp = strings.TrimSpace(statusUp)

		if statusUp == "" {
			return pError.New("statusUp can't be empty")
		}

		cfg.StatusUp = statusUp
		return nil
	})
}

func StatusCodeDown(statusCodeDown int) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		if statusCodeDown <= 0 {
			return pError.New("statusCodeDown must be positive")
		}

		cfg.StatusCodeDown = statusCodeDown
		return nil
	})
}

func StatusDown(statusDown string) options.Option[Config] {
	return options.OptionFunc[Config](func(cfg *Config) error {
		statusDown = strings.TrimSpace(statusDown)

		if statusDown == "" {
			return pError.New("statusDown can't be empty")
		}

		cfg.StatusDown = statusDown
		return nil
	})
}
