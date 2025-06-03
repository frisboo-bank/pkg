package factory

import (
	"frisboo-bank/pkg/di"
	"frisboo-bank/pkg/di/config"
	"frisboo-bank/pkg/di/simpledi"
)

func NewInstance(cfg *config.DiOptions) di.Container {
	return simpledi.NewSimpleDi(cfg)
}
