package dig

import (
	"fmt"

	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"go.uber.org/dig"
)

func (a *digAdapter) RegisterDecorators(decorators ...decorator.Decorator) error {
	for id, d := range decorators {
		name := fmt.Sprintf("decorator-%d", id)
		if err := a.RegisterDecorator(name, d); err != nil {
			return err
		}
	}
	return nil
}

func (a *digAdapter) RegisterDecorator(name string, d decorator.Decorator) error {
	cfg := decorator.Config{}
	if err := options.Apply(&cfg, d.Options()...); err != nil {
		return syserrors.Wrap(err, "failed to apply decorator %s options", name)
	}

	fn, err := wrapFuncWithDigIn(d.Fn(), cfg.NamedDeps, name)
	if err != nil {
		return syserrors.Wrap(err, "failed to adapt decorator %s named deps", name)
	}

	opts := toDigDecoratorOptions(cfg)

	if err := a.dig.Decorate(fn, opts...); err != nil {
		return syserrors.Wrap(err, "failed to register decorator %s", name)
	}
	return nil
}

func toDigDecoratorOptions(cfg decorator.Config) []dig.DecorateOption {
	var result []dig.DecorateOption

	if bcb, ok := any(cfg.BeforeCallback).(dig.BeforeCallback); ok && bcb != nil {
		result = append(result, dig.WithDecoratorBeforeCallback(bcb))
	}

	if cb, ok := any(cfg.Callback).(dig.Callback); ok && cb != nil {
		result = append(result, dig.WithDecoratorCallback(cb))
	}

	if info, ok := cfg.Info.(*dig.DecorateInfo); ok && info != nil {
		result = append(result, dig.FillDecorateInfo(info))
	}

	return result
}
