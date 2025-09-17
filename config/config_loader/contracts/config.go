package contracts

import (
	"frisboo-bank/pkg/environment"
)

type ConfigLoader interface {
	Load(env environment.Environment, cfg any) error
	LoadKey(env environment.Environment, cfg any, key string) error
	HasKey(env environment.Environment, key string) (bool, error)
}
