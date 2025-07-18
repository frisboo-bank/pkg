package contracts

import (
	"frisboo-bank/pkg/environment"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type ConfigLoader interface {
	WithViper(viper *viper.Viper) ConfigLoader
	WithDecodeHooks(decodeHooks ...mapstructure.DecodeHookFunc) ConfigLoader
	WithConfigPath(path string) ConfigLoader
	WithConfigName(name string) ConfigLoader

	LoadConfig(env environment.Environment, cfg any) error
	LoadConfigByKey(key string, env environment.Environment, cfg any) error
}
