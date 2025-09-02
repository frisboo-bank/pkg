package config

import (
	"strings"

	"frisboo-bank/pkg/config"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
)

var _ config.Validatable = (*Config)(nil)

type Config struct {
	ConfigPath      string                        `mapstructure:"configPath"`
	ConfigName      string                        `mapstructure:"configName"`
	Debug           bool                          `mapstructure:"debug"`
	EnvKeyReplacer  map[string]string             `mapstructure:"envKeyReplacer"`
	EnvPrefix       string                        `mapstructure:"envPrefix"`
	Viper           *viper.Viper                  `mapstructure:"-"`
	DecodeHookFuncs []mapstructure.DecodeHookFunc `mapstructure:"-"`
}

func Default() *Config {
	return &Config{
		ConfigPath: "./configs",
		ConfigName: "application",
		Debug:      false,
		EnvPrefix:  "APP_",
	}
}

func (c *Config) Validate() error {
	var errs *multierror.Error

	if strings.TrimSpace(c.ConfigName) == "" {
		errs = multierror.Append(errs, syserrors.CantBeEmptyError("ConfigName"))
	}

	return errs.ErrorOrNil()
}

func New(opts ...Option) (*Config, error) {
	return options.New(Default, opts...)
}
