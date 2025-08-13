package config

import (
	"frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/environment"
)

func LoadConfig[T any](loader contracts.ConfigLoader, env environment.Environment, key string) (*T, error) {
	cfg := new(T)

	err := loader.LoadConfigByKey(key, env, cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
