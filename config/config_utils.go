package config

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"frisboo-bank/pkg/constants"
	customErrors "frisboo-bank/pkg/custom_errors"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection/typemapper"
	"frisboo-bank/pkg/utils"

	envParser "github.com/caarlos0/env/v11"
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

var (
	ErrReadConfig      = errors.New("config: failed to read config")
	ErrUnmarchalConfig = errors.New("config: failed to unmarshal config")
)

func BindConfig[T any](env environment.Environment) (T, error) {
	return BindConfigKey[T]("", env)
}

func BindConfigKey[T any](key string, env environment.Environment) (T, error) {
	cfg := typemapper.GenericInstanceByT[T]()

	defaults.SetDefaults(cfg)

	configPath, err := getConfigPath()
	if err != nil {
		return *new(T), err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("application")

	if err := viper.ReadInConfig(); err != nil {
		return *new(T), customErrors.WrapWith(ErrReadConfig, err)
	}

	viper.SetConfigName(fmt.Sprintf("application.%s", string(env)))

	err = viper.MergeInConfig()

	targetError := &viper.ConfigFileNotFoundError{}
	if errors.As(err, &targetError) {
		return *new(T), customErrors.WrapWith(ErrReadConfig, err)
	}

	if len(key) == 0 {
		err = viper.Unmarshal(cfg)
	} else {
		err = viper.UnmarshalKey(key, cfg)
	}

	if err != nil {
		return *new(T), customErrors.WrapWith(ErrUnmarchalConfig, err)
	}

	viper.AutomaticEnv()

	if err := envParser.Parse(cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return cfg, nil
}

func getConfigPath() (string, error) {
	configPath := viper.GetString(constants.CONFIG_PATH_ENV)

	if strings.TrimSpace(configPath) != "" {
		return configPath, nil
	}

	appRootPath, err := generateConfigRootPath()
	if err != nil {
		return "", err
	}

	return appRootPath, nil
}

func generateConfigRootPath() (string, error) {
	rootPath := viper.GetString(constants.APP_ROOT_PATH_ENV)

	if strings.TrimSpace(rootPath) != "" {
		return rootPath, nil
	}

	rootPath, err := utils.GetProjectRootWorkingDirectory()
	if err != nil {
		return "", err
	}

	return path.Join(rootPath, "configs"), nil
}
