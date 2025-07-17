package config

import (
	"fmt"
	"frisboo-bank/pkg/config/contracts"
	"frisboo-bank/pkg/config/options"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/utils"
	"path"
	"strings"

	envparse "github.com/caarlos0/env/v11"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type configLoader struct {
	viper       *viper.Viper
	decodeHooks []mapstructure.DecodeHookFunc
	configPath  string
	configName  string
}

var _ contracts.ConfigLoader = (*configLoader)(nil)

func (c *configLoader) WithViper(viper *viper.Viper) contracts.ConfigLoader {
	c.viper = viper
	return c
}

func (c *configLoader) WithDecodeHooks(decodeHooks ...mapstructure.DecodeHookFunc) contracts.ConfigLoader {
	c.decodeHooks = decodeHooks
	return c
}

func (c *configLoader) WithConfigPath(path string) contracts.ConfigLoader {
	c.configPath = path
	return c
}

func (c *configLoader) WithConfigName(name string) contracts.ConfigLoader {
	c.configName = name
	return c
}

func NewConfigLoader() contracts.ConfigLoader {
	return &configLoader{
		viper:      viper.New(),
		configPath: options.ConfigPath,
		configName: options.ConfigName,
	}
}

func (c *configLoader) LoadConfig(env environment.Environment, cfg any) error {
	return c.LoadConfigByKey("", env, cfg)
}

func (c *configLoader) LoadConfigByKey(key string, env environment.Environment, cfg any) error {
	configPath := c.getConfigPath()
	if configPath == "" {
		return fmt.Errorf("you haven't set a path where to look for config files, your can set the Env: %s or %s",
			constants.CONFIG_PATH_ENV,
			constants.APP_ROOT_PATH_ENV,
		)
	}

	if !path.IsAbs(configPath) {
		rootPath, err := utils.GetProjectRootWorkingDirectory()
		if err != nil {
			return err
		}

		configPath = path.Join(rootPath, configPath)
	}

	configName := fmt.Sprintf("%s.%s", c.configName, string(env))

	c.viper.AddConfigPath(configPath)
	c.viper.SetConfigName(configName)

	if err := c.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config with error: %w", err)
	}

	decoderHooks := []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(c.decodeHooks...)),
	}

	var err error
	if len(key) == 0 {
		err = c.viper.Unmarshal(&cfg, decoderHooks...)
	} else {
		subKey := c.viper.Sub(key)
		if subKey == nil {
			return fmt.Errorf("config-loader: failed to unmarshal config with error: no such key %s", key)
		}

		err = c.viper.UnmarshalKey(key, cfg, decoderHooks...)
	}

	if err != nil {
		return fmt.Errorf("config-loader: failed to unmarshal config with error: %w", err)
	}

	c.viper.AutomaticEnv()

	if err := envparse.Parse(cfg); err != nil {
		fmt.Printf("config-loader: %+v\n", err)
	}

	return nil
}

func (cl *configLoader) getConfigPath() string {
	configPath := strings.TrimSpace(cl.viper.GetString(constants.CONFIG_PATH_ENV))
	if configPath != "" {
		return configPath
	}

	configPath = strings.TrimSpace(cl.viper.GetString(constants.APP_ROOT_PATH_ENV))
	if configPath != "" {
		return configPath
	}

	return strings.TrimSpace(cl.configPath)
}
