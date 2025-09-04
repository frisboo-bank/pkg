package configloader

import (
	"fmt"
	"path"
	"reflect"
	"strings"

	"frisboo-bank/pkg/config/config_loader/config"
	"frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/utils"
	"frisboo-bank/pkg/validation"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type configLoader struct {
	cfg   *config.Config
	viper *viper.Viper
}

func New(cfg *config.Config) contracts.ConfigLoader {
	validation.Assert(cfg != nil, syserrors.CantBeNilError("cfg"))

	vi := cfg.Viper
	if vi == nil {
		vi = viper.New()
	}

	vi.AutomaticEnv()

	if cfg.Debug {
		vi.Debug()
	}

	if cfg.EnvPrefix != "" {
		vi.SetEnvPrefix(cfg.EnvPrefix)
	}

	if cfg.EnvKeyReplacer != nil {
		replacer := make([]string, len(cfg.EnvKeyReplacer)*2)
		for match, replace := range cfg.EnvKeyReplacer {
			replacer = append(replacer, match, replace)
		}
		vi.SetEnvKeyReplacer(strings.NewReplacer(replacer...))
	}

	return &configLoader{
		cfg,
		vi,
	}
}

func (c *configLoader) Load(env environment.Environment, cfg any) error {
	return c.loadByKey("", env, cfg)
}

func (c *configLoader) LoadByKey(key string, env environment.Environment, cfg any) error {
	return c.loadByKey(key, env, cfg)
}

func (c *configLoader) loadByKey(key string, env environment.Environment, cfg any) error {
	cfgType := reflect.ValueOf(cfg).Kind()
	if cfgType != reflect.Pointer {
		return syserrors.Newf("cfg must be a pointer to a struct: currently got %s", cfgType.String())
	}

	configPath := c.cfg.ConfigPath
	if configPath == "" {
		return syserrors.Newf("you haven't set a path where to look for config files, your can set the Env: %s or %s",
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

	configName := fmt.Sprintf("%s.%s", c.cfg.ConfigName, string(env))

	if c.cfg.Debug {
		fmt.Printf("try to load the config file with name %s in the path %s\n", configName, configPath)
	}

	c.viper.SetConfigName(configName)
	c.viper.AddConfigPath(configPath)

	if err := c.viper.ReadInConfig(); err != nil {
		return syserrors.Newf("failed to read config with error: %w", err)
	}

	decoderHooks := []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(c.cfg.DecodeHookFuncs...)),
	}

	var err error
	if len(key) == 0 {
		err = c.viper.Unmarshal(&cfg, decoderHooks...)
	} else {
		subKey := c.viper.Sub(key)
		if subKey == nil {
			return nil
		}
		err = c.viper.UnmarshalKey(key, &cfg, decoderHooks...)
	}

	if err != nil {
		return syserrors.Newf("failed to unmarshal config with error: %w", err)
	}

	return nil
}
