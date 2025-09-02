package config

import (
	"strings"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Option = options.OptionFn[Config]

var ConfigPath = options.OptionErr(func(c *Config, configPath string) error {
	configPath = strings.TrimSpace(configPath)
	if configPath == "" {
		return syserrors.CantBeEmptyError("ConfigPath")
	}
	c.ConfigPath = configPath
	return nil
})

var ConfigName = options.OptionErr(func(c *Config, configName string) error {
	configName = strings.TrimSpace(configName)
	if configName == "" {
		return syserrors.CantBeEmptyError("ConfigName")
	}
	c.ConfigName = configName
	return nil
})

var EnvKeyReplacer = options.OptionErr(func(c *Config, envKeyReplacer map[string]string) error {
	if len(envKeyReplacer) == 0 {
		return syserrors.CantBeEmptyError("EnvKeyReplacer")
	}
	c.EnvKeyReplacer = envKeyReplacer
	return nil
})

var EnvPrefix = options.Option(func(c *Config, envPrefix string) {
	c.EnvPrefix = strings.TrimSpace(envPrefix)
})

var Debug = options.Option(func(c *Config, debug bool) {
	c.Debug = debug
})

var Viper = options.OptionErr(func(c *Config, viper *viper.Viper) error {
	if viper == nil {
		return syserrors.CantBeNilError("Viper")
	}
	c.Viper = viper
	return nil
})

var DecodeHookFuncs = options.VarOptionErr(func(c *Config, decodeHookFuncs ...mapstructure.DecodeHookFunc) error {
	if len(decodeHookFuncs) == 0 {
		return syserrors.CantBeEmptyError("DecodeHookFuncs")
	}
	c.DecodeHookFuncs = append(c.DecodeHookFuncs, decodeHookFuncs...)
	return nil
})
