package configloader

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"frisboo-bank/pkg/config/config_loader/config"
	"frisboo-bank/pkg/config/config_loader/contracts"
	"frisboo-bank/pkg/constants"
	"frisboo-bank/pkg/environment"
	"frisboo-bank/pkg/reflection"
	"frisboo-bank/pkg/syserrors"
	"frisboo-bank/pkg/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

var _ contracts.ConfigLoader = (*configLoader)(nil)

type configLoader struct {
	cfg   config.Config
	viper *viper.Viper
	mu    sync.RWMutex
}

func New(opts ...config.Option) (contracts.ConfigLoader, error) {
	cfg, err := config.New(opts...)
	if err != nil {
		return nil, syserrors.Wrap(err, "failed to instantiate config")
	}

	vi := cfg.Viper
	if vi == nil {
		vi = viper.New()
	}

	if cfg.Debug {
		vi.Debug()
	}

	vi.AutomaticEnv()

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
		cfg:   cfg,
		viper: vi,
	}, nil
}

func (c *configLoader) Load(env environment.Environment, cfg any) error {
	if err := c.loadConfigFile(env); err != nil {
		return err
	}
	return c.unmarshal("", cfg)
}

func (c *configLoader) LoadKey(env environment.Environment, cfg any, key string) error {
	if err := c.loadConfigFile(env); err != nil {
		return err
	}
	return c.unmarshal(key, cfg)
}

func (c *configLoader) LoadComposableKey(env environment.Environment, cfg any, keys ...string) error {
	if len(keys) == 0 {
		return syserrors.CantBeEmptyError("keys")
	}
	if err := c.loadConfigFile(env); err != nil {
		return err
	}
	for _, k := range keys {
		if err := c.unmarshal(k, cfg); err != nil {
			return err
		}
	}
	return nil
}

func (c *configLoader) HasKey(env environment.Environment, key string) (bool, error) {
	if key == "" {
		return false, syserrors.CantBeEmptyError("key")
	}
	if err := c.loadConfigFile(env); err != nil {
		return false, err
	}
	return c.viper.Sub(key) != nil, nil
}

func (c *configLoader) loadConfigFile(env environment.Environment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := validation.Validate(env, validation.Required); err != nil {
		return err
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

	configName := fmt.Sprintf("%s.%s", c.cfg.ConfigName, env)

	if c.cfg.Debug {
		fmt.Printf("try to load the config file with name %s in the path %s\n", configName, configPath)
	}

	c.viper.SetConfigName(configName)
	c.viper.AddConfigPath(configPath)

	if err := c.viper.ReadInConfig(); err != nil {
		return syserrors.Wrapf(err, "failed to read config %s", configName)
	}

	return nil
}

func (c *configLoader) unmarshal(key string, target any) error {
	if !reflection.IsPointer(target) {
		return syserrors.Newf("target must be a pointer: got %s", reflection.GetKind(target))
	}

	opts := []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(c.cfg.DecodeHookFuncs...)),
	}

	if key == "" {
		return c.viper.Unmarshal(target, opts...)
	}

	if c.viper.Sub(key) == nil {
		return nil
	}
	return c.viper.UnmarshalKey(key, target, opts...)
}
