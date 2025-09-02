package config

import "frisboo-bank/pkg/config"

var _ config.Validatable = (*Config)(nil)

type Config struct{}

func Default() *Config {
	return &Config{}
}

func (c *Config) Validate() error {
	return nil
}
