package dig

import (
	"fmt"

	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"

	"go.uber.org/dig"
)

func (a *digAdapter) RegisterInvokers(invokers ...invoker.Invoker) error {
	for id, i := range invokers {
		name := fmt.Sprintf("invoker-%d", id)
		if err := a.RegisterInvoker(name, i); err != nil {
			return err
		}
	}
	return nil
}

func (a *digAdapter) RegisterInvoker(name string, i invoker.Invoker) error {
	cfg := invoker.Config{}
	if err := options.Apply(&cfg, i.Options()...); err != nil {
		return syserrors.Wrapf(err, "failed to apply invoker %s options", name)
	}

	fn := i.Fn()
	if len(cfg.NamedDeps) > 0 {
		var err error
		fn, err = wrapFuncWithNamedInputs(fn, cfg.NamedDeps, name)
		if err != nil {
			return syserrors.Wrapf(err, "failed to adapt invoker %s named deps", name)
		}
	}

	opts := toDigInvokerOptions(cfg)

	if err := a.dig.Invoke(fn, opts...); err != nil {
		return syserrors.Wrapf(err, "failed to register invoker %s", name)
	}

	return nil
}

func toDigInvokerOptions(cfg invoker.Config) []dig.InvokeOption {
	var result []dig.InvokeOption

	if info, ok := cfg.Info.(*dig.InvokeInfo); ok && info != nil {
		result = append(result, dig.FillInvokeInfo(info))
	}

	return result
}
