package config

import "frisboo-bank/pkg/options"

type Config struct{}

var defaultConfig = &Config{}

func Apply() *options.OptionBuilder[Config] {
	return options.Apply(defaultConfig)
}
