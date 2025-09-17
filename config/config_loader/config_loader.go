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

	validation "github.com/go-ozzo/ozzo-validation/v4"
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

/**
 *
 * Commands
 *
 */
func (c *configLoader) Load(env environment.Environment, target any) error {
	if !reflection.IsPointer(target) {
		return syserrors.Newf("target must be a pointer: got %s", reflection.GetKind(target))
	}
	if err := c.ensureLoaded(env); err != nil {
		return err
	}
	return c.unmarshal("", target)
}

func (c *configLoader) LoadKey(env environment.Environment, target any, key string) error {
	if !reflection.IsPointer(target) {
		return syserrors.Newf("target must be a pointer: got %s", reflection.GetKind(target))
	}
	if err := c.ensureLoaded(env); err != nil {
		return err
	}
	return c.unmarshal(key, target)
}

func (c *configLoader) HasKey(env environment.Environment, key string) (bool, error) {
	if key == "" {
		return false, syserrors.CantBeEmptyError("key")
	}
	if err := c.ensureLoaded(env); err != nil {
		return false, err
	}
	return c.keyExists(key), nil
}

/**
 *
 * Internal Helpers
 *
 */
func (c *configLoader) ensureLoaded(env environment.Environment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := validation.Validate(env, validation.Required); err != nil {
		return syserrors.Wrap(err, "env validation failed")
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
			return syserrors.Wrap(err, "failed to resolve project root")
		}
		configPath = path.Join(rootPath, configPath)
	}

	configName := fmt.Sprintf("%s.%s", c.cfg.ConfigName, env)
	c.viper.SetConfigName(configName)
	c.viper.AddConfigPath(configPath)

	if err := c.viper.ReadInConfig(); err != nil {
		return syserrors.Wrapf(err, "failed to read config %s in path %s", configName, configPath)
	}

	if c.cfg.Debug {
		if fileUsed := c.viper.ConfigFileUsed(); fileUsed != "" {
			fmt.Printf("loaded config file %s\n", fileUsed)
		}
	}

	return nil
}

func (c *configLoader) unmarshal(key string, target any) error {
	opts := []viper.DecoderConfigOption{
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(c.cfg.DecodeHookFuncs...)),
	}

	if key == "" {
		if err := c.viper.UnmarshalExact(target, opts...); err != nil {
			return syserrors.Wrap(err, "failed to unmarshal root")
		}
		return nil
	}

	if !c.keyExists(key) {
		return syserrors.Newf("required key %s not found", key)
	}

	sub := c.viper.Sub(key)
	if sub == nil {
		return syserrors.Newf("required key %s not found", key)
	}
	if err := sub.UnmarshalExact(target, opts...); err != nil {
		return syserrors.Wrapf(err, "failed to unmarshal key %s", key)
	}
	return nil
}

func (c *configLoader) keyExists(key string) bool {
	if c.viper.IsSet(key) {
		return true
	}
	sub := c.viper.Sub(key)
	return sub != nil && len(sub.AllSettings()) > 0
}
