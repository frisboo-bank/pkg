package config

import (
	"time"

	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	validationIs "github.com/go-ozzo/ozzo-validation/v4/is"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	DB           string        `mapstructure:"db"`
	Username     string        `mapstructure:"username"`
	Password     string        `mapstructure:"password"`
	PoolSize     int           `mapstructure:"poolSize"`
	DialTimeout  time.Duration `mapstructure:"dialTimeout"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
}

func Default() *Config {
	return &Config{
		DialTimeout:  0,
		PoolSize:     0,
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
}

type Option = options.OptionFn[Config]

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required, validationIs.Host),
		validation.Field(&c.Port, validation.Required, validationIs.Port),
		validation.Field(&c.DB, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.PoolSize, validation.Required, validation.Min(0)),
		validation.Field(&c.DialTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.ReadTimeout, validation.Required, validation.Min(0)),
		validation.Field(&c.WriteTimeout, validation.Required, validation.Min(0)),
	)
}
