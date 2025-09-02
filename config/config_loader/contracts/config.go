package contracts

import (
	"frisboo-bank/pkg/environment"
)

type ConfigLoader interface {
	Load(env environment.Environment, cfg any) error
	LoadByKey(key string, env environment.Environment, cfg any) error
}
