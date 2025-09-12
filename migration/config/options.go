package config

import (
	"strings"

	"frisboo-bank/pkg/options"
)

type Option = options.OptionFn[Config]

var DB = options.Option(func(c *Config, dbKey string) {
	c.DBKey = strings.TrimSpace(dbKey)
})

var MigrationsDir = options.Option(func(c *Config, migrationsDir string) {
	c.MigrationsDir = strings.TrimSpace(migrationsDir)
})
