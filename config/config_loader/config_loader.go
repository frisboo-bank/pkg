package configloader

import (
	"errors"
	"fmt"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/utils"
	"path"
	"strings"

	envParser "github.com/caarlos0/env/v11"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type ConfigBinderConfig struct {
	Viper       *viper.Viper
	DecodeHooks []mapstructure.DecodeHookFunc
}

type ConfigBinder struct {
	viper      *viper.Viper
	decodeHook mapstructure.DecodeHookFunc
}

func New(options ConfigBinder) *ConfigBinder {
	if options.viper == nil {
		options.viper = viper.New()
	}

	return &ConfigBinder{
		viper:      options.viper,
		decodeHook: options.decodeHook,
	}
}

// BindConfig loads the config struct from Viper using the binder's decode hook.
func (cb *ConfigBinder) BindConfig(env environment.Environment, cfg any) error {
	return cb.BindConfigKey("", env, cfg)
}

// BindConfigKey loads a config struct from a specific key
func (cb *ConfigBinder) BindConfigKey(key string, env environment.Environment, cfg any) error {
	configPath, err := getConfigPath(cb.viper)
	if err != nil {
		return err
	}

	cb.viper.AddConfigPath(configPath)
	cb.viper.SetConfigName("application")

	if err := cb.viper.ReadInConfig(); err != nil {
		return fmt.Errorf("config-binder: failed to read config with error: %w", err)
	}

	cb.viper.SetConfigName(fmt.Sprintf("application.%s", string(env)))

	err = cb.viper.MergeInConfig()

	targetError := &viper.ConfigFileNotFoundError{}
	if errors.As(err, &targetError) {
		return fmt.Errorf("config-binder: failed to read config with error: %w", err)
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: cb.decodeHook,
		Result:     cfg,
		TagName:    "mapstructure",
	})
	if err != nil {
		return fmt.Errorf("config-binder: failed to unmarshal config with error: %w", err)
	}

	if len(key) == 0 {
		err = decoder.Decode(cb.viper.AllSettings())
	} else {
		subKey := cb.viper.Sub(key)
		if subKey == nil {
			return fmt.Errorf("config-binder: failed to unmarshal config with error: no such key %s", key)
		}

		err = decoder.Decode(subKey.AllSettings())
	}

	if err != nil {
		return fmt.Errorf("config-binder: failed to unmarshal config with error: %w", err)
	}

	cb.viper.AutomaticEnv()

	if err := envParser.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return nil
}

func getConfigPath(v *viper.Viper) (string, error) {
	configPath := v.GetString(constants.CONFIG_PATH_ENV)

	if strings.TrimSpace(configPath) != "" {
		return configPath, nil
	}

	appRootPath, err := generateConfigRootPath(v)
	if err != nil {
		return "", err
	}

	return appRootPath, nil
}

func generateConfigRootPath(v *viper.Viper) (string, error) {
	rootPath := v.GetString(constants.APP_ROOT_PATH_ENV)
	if strings.TrimSpace(rootPath) != "" {
		return rootPath, nil
	}

	rootPath, err := utils.GetProjectRootWorkingDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(rootPath, "configs"), nil
}
