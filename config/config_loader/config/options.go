package config

import (
	"strings"

	"frisboo-bank/pkg/options"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Option = options.OptionFn[Config]

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
