package config

import (
	configloaderContracts "frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/options"
)

type Config struct {
	DBKey         string `mapstructure:"db"`
	MigrationsDir string `mapstructure:"migrationsDir"`
}

func Default() *Config {
	return &Config{
		MigrationsDir: "./migrations",
	}
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}

func Load(loader configloaderContracts.ConfigLoader, env environment.Environment, opts ...Option) (*Config, error) {
	cfg := Default()
	if err := loader.LoadByKey("migration", env, cfg); err != nil {
		return nil, err
	}
	return options.New(func() *Config { return cfg }, opts...)
}
