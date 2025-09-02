package config

import (
	"strings"

	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
)

type Option options.OptionFn[Config]

var DB = options.OptionErr(func(c *Config, db string) error {
	db = strings.TrimSpace(db)
	if db == "" {
		return syserrors.CantBeEmptyError("DB")
	}
	c.DBKey = db
	return nil
})

var MigrationsDir = options.OptionErr(func(c *Config, migrationsDir string) error {
	migrationsDir = strings.TrimSpace(migrationsDir)
	if migrationsDir == "" {
		return syserrors.CantBeEmptyError("MigrationsDir")
	}
	c.MigrationsDir = migrationsDir
	return nil
})
