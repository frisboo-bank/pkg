package dig

import (
	"reflect"

	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/options"
	"frisboo-bank/pkg/syserrors"
	waiterContracts "frisboo-bank/pkg/waiter/contracts"

	"go.uber.org/dig"
)

func (a *digAdapter) RegisterHooks(hooks ...hook.Hooks) error {
	for _, h := range hooks {
		if err := a.RegisterHook(h); err != nil {
			return err
		}
	}
	return nil
}

func (a *digAdapter) RegisterHook(h hook.Hooks) error {
	name := h.Name()

	cfg := hook.Config{}
	if err := options.Apply(&cfg, h.Options()...); err != nil {
		return syserrors.Wrapf(err, "failed to apply hook %s options", name)
	}

	startCtor := h.StartConstructor()
	stopCtor := h.StopConstructor()

	if startCtor == nil && stopCtor == nil {
		return nil
	}

	var err error
	if startCtor != nil {
		startCtor, err = wrapHookConstructorWithNamedDeps(startCtor, cfg.NamedDeps, name, "start")
		if err != nil {
			return err
		}
	}
	if stopCtor != nil {
		stopCtor, err = wrapHookConstructorWithNamedDeps(stopCtor, cfg.NamedDeps, name, "stop")
		if err != nil {
			return err
		}
	}

	opts := toDigHookOptions(cfg)

	startGroup := name + "-start"
	if startCtor != nil {
		if err := a.dig.Provide(startCtor, append(opts, dig.Group(startGroup))...); err != nil {
			return syserrors.Wrapf(err, "failed to register hook start %s", startGroup)
		}
	}
	stopGroup := name + "-stop"
	if stopCtor != nil {
		if err := a.dig.Provide(stopCtor, append(opts, dig.Group(stopGroup))...); err != nil {
			return syserrors.Wrapf(err, "failed to register hook stop %s", stopGroup)
		}
	}

	a.hooks = append(a.hooks, name)

	return nil
}

func (a *digAdapter) resolveHooks() ([]waiterContracts.WaiterHook, error) {
	hooks := make([]waiterContracts.WaiterHook, len(a.hooks))
	for i, h := range a.hooks {
		h, err := a.resolveHook(h)
		if err != nil {
			return nil, err
		}
		hooks[i] = h
	}
	return hooks, nil
}

func (a *digAdapter) resolveHook(hook string) (waiterContracts.WaiterHook, error) {
	wh := waiterContracts.WaiterHook{Name: hook}

	startGroup := hook + "-start"
	startFns, err := resolveDynamicGroup[[]waiterContracts.WaitFunc](
		a.dig,
		startGroup,
		reflect.TypeOf([]waiterContracts.WaitFunc{}),
	)
	if err != nil {
		return wh, syserrors.Wrapf(err, "failed to resolve start hook %s", startGroup)
	}
	if len(startFns) > 0 {
		wh.Wait = startFns[0]
	}

	stopGroup := hook + "-stop"
	stopFns, err := resolveDynamicGroup[[]waiterContracts.CleanupFunc](
		a.dig,
		stopGroup,
		reflect.TypeOf([]waiterContracts.CleanupFunc{}),
	)
	if err != nil {
		return wh, syserrors.Wrapf(err, "failed to resolve stop hook %s", stopGroup)
	}
	if len(stopFns) > 0 {
		wh.Cleanup = stopFns[0]
	}

	return wh, nil
}

func wrapHookConstructorWithNamedDeps(ctor any, namedDeps map[string]string, hookName, phase string) (any, error) {
	if ctor == nil || len(namedDeps) == 0 {
		return ctor, nil
	}

	origVal := reflect.ValueOf(ctor)
	origType := origVal.Type()

	if origType.Kind() != reflect.Func {
		return nil, syserrors.Newf("hook %s (%s): constructor not a func", hookName, phase)
	}
	if origType.NumIn() != 1 {
		return ctor, nil
	}

	paramType := origType.In(0)
	ptrParam := false
	if paramType.Kind() == reflect.Pointer {
		ptrParam = true
		paramType = paramType.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		// Not the supported shape; leave untouched.
		return ctor, nil
	}

	// Build dig.In wrapper struct fields: dig.In + each exported tagged field.
	fields := []reflect.StructField{
		{
			Name:      "In",
			Type:      reflect.TypeOf(dig.In{}),
			Anonymous: true,
		},
	}

	// Track which original fields we will rehydrate (index mapping).
	var originalFieldIndexes []int
	for i := 0; i < paramType.NumField(); i++ {
		sf := paramType.Field(i)
		if !sf.IsExported() {
			continue
		}
		tagName, has := sf.Tag.Lookup("name")
		if !has || tagName == "" {
			continue
		}
		actual := namedDeps[tagName]
		if actual == "" {
			// If user didn't provide a mapping, treat symbolic as actual.
			actual = tagName
		}
		fields = append(fields, reflect.StructField{
			Name: sf.Name,
			Type: sf.Type,
			Tag:  reflect.StructTag(`name:"` + actual + `"`),
		})
		originalFieldIndexes = append(originalFieldIndexes, i)
	}

	// If we didn't find any tagged fields, no adaptation needed.
	if len(originalFieldIndexes) == 0 {
		return ctor, nil
	}

	inStruct := reflect.StructOf(fields)

	// Mirror outputs directly.
	numOut := origType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for i := range numOut {
		outTypes[i] = origType.Out(i)
	}

	wrapperType := reflect.FuncOf([]reflect.Type{inStruct}, outTypes, false)

	wrapped := reflect.MakeFunc(wrapperType, func(args []reflect.Value) []reflect.Value {
		inVal := args[0]

		// Recreate the original param (struct or pointer to struct).
		var param reflect.Value
		if ptrParam {
			param = reflect.New(paramType)
		} else {
			param = reflect.New(paramType).Elem()
		}

		for idx, origFieldIdx := range originalFieldIndexes {
			fieldVal := inVal.Field(idx + 1) // +1 skip embedded dig.In
			if ptrParam {
				param.Elem().Field(origFieldIdx).Set(fieldVal)
			} else {
				param.Field(origFieldIdx).Set(fieldVal)
			}
		}

		callArgs := []reflect.Value{param}
		return origVal.Call(callArgs)
	})

	return wrapped.Interface(), nil
}

func toDigHookOptions(cfg hook.Config) []dig.ProvideOption {
	var result []dig.ProvideOption

	if len(cfg.As) > 0 {
		result = append(result, dig.As(cfg.As...))
	}

	if cfg.Export {
		result = append(result, dig.Export(true))
	}

	if cfg.Group != "" {
		result = append(result, dig.Group(cfg.Group))
	}

	if cfg.LocationPC != 0 {
		result = append(result, dig.LocationForPC(cfg.LocationPC))
	}

	if cfg.Name != "" {
		result = append(result, dig.Name(cfg.Name))
	}

	return result
}
