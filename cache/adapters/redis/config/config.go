package config

import (
	"time"

	"frisboo-bank/pkg/config"

	"github.com/hashicorp/go-multierror"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	Addr         string        `mapstructure:"addr"`
	DB           int           `mapstructure:"db"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	Host         string        `mapstructure:"host"`
	Password     string        `mapstructure:"password"`
	PoolSize     int           `mapstructure:"poolSize"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	Username     string        `mapstructure:"username"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

func Default() *Config {
	return &Config{
		Addr:         "",
		DB:           0,
		DialTimeout:  0,
		Host:         "",
		Password:     "",
		PoolSize:     0,
		Port:         0,
		ReadTimeout:  0,
		Username:     "",
		WriteTimeout: 0,
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	return errs.ErrorOrNil()
}
