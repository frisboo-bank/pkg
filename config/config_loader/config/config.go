package config

import (
	"frisboo-bank/pkg/options"
	cValidation "frisboo-bank/pkg/validation"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var _ cValidation.Validatable = (*Config)(nil)

type Config struct {
	ConfigName      string                        `mapstructure:"configName"`
	ConfigPath      string                        `mapstructure:"configPath"`
	Debug           bool                          `mapstructure:"debug"`
	DecodeHookFuncs []mapstructure.DecodeHookFunc `mapstructure:"-"`
	EnvKeyReplacer  map[string]string             `mapstructure:"envKeyReplacer"`
	EnvPrefix       string                        `mapstructure:"envPrefix"`
	Viper           *viper.Viper                  `mapstructure:"-"`
}

func Default() Config {
	return Config{
		ConfigName: "application",
		ConfigPath: "./configs",
		Debug:      false,
		EnvPrefix:  "APP_",
	}
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(c.ConfigName, validation.Required),
		validation.Field(c.ConfigPath, validation.Required),
		validation.Field(c.EnvPrefix, validation.Required),
	)
}

func New(opts ...Option) (Config, error) {
	var zero Config

	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return zero, err
	}

	return base, nil
}
