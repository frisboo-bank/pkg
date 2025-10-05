package config

import (
	"strings"

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

type Option = options.OptionFn[Config]

func New(opts ...Option) (Config, error) {
	var zero Config
	base := Default()
	if err := options.Apply(&base, opts...); err != nil {
		return zero, err
	}
	return base, nil
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ConfigName, validation.Required),
		validation.Field(&c.ConfigPath, validation.Required),
		validation.Field(&c.EnvPrefix, validation.Required),
	)
}

var ConfigPath = options.Option(func(c *Config, configPath string) {
	c.ConfigPath = strings.TrimSpace(configPath)
})

var ConfigName = options.Option(func(c *Config, configName string) {
	c.ConfigName = strings.TrimSpace(configName)
})

var EnvKeyReplacer = options.Option(func(c *Config, envKeyReplacer map[string]string) {
	c.EnvKeyReplacer = envKeyReplacer
})

var EnvPrefix = options.Option(func(c *Config, envPrefix string) {
	c.EnvPrefix = strings.TrimSpace(envPrefix)
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Viper = options.Option(func(c *Config, viper *viper.Viper) {
	c.Viper = viper
})

var DecodeHookFuncs = options.VarOption(func(c *Config, decodeHookFuncs ...mapstructure.DecodeHookFunc) {
	c.DecodeHookFuncs = append(c.DecodeHookFuncs, decodeHookFuncs...)
})
