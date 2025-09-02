package cache

import (
	"frisboo-bank/pkg/container/dependencies/module"
)

var ModuleFunc = func() module.Module {
	m := module.ModuleFunc("cache")

	return m
}
