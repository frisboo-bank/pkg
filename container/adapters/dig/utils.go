package dig

import (
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"
	"frisboo-bank/pkg/syserrors"

	"go.uber.org/dig"
)

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

func toDigInvokerOptions(cfg invoker.Config) []dig.InvokeOption {
	var result []dig.InvokeOption

	if info, ok := cfg.Info.(*dig.InvokeInfo); ok && info != nil {
		result = append(result, dig.FillInvokeInfo(info))
	}

	return result
}

func toDigProvideOptions(cfg provider.Config) []dig.ProvideOption {
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

// resolveDynamicGroup dynamically resolves a dig group into a value of type T.
// It constructs a struct input with a group tag and invokes the container.
// It returns the resolved group or an error if resolution fails.
func resolveDynamicGroup[T any](
	container *dig.Container,
	groupName string,
	sliceType reflect.Type,
) (T, error) {
	groupTag := fmt.Sprintf(`group:"%s"`, groupName)

	inputType := reflect.StructOf([]reflect.StructField{
		{
			Name:      "In",
			Type:      reflect.TypeOf(dig.In{}),
			Anonymous: true,
		},
		{
			Name: "Funcs",
			Type: sliceType,
			Tag:  reflect.StructTag(groupTag),
		},
	})

	input := reflect.New(inputType).Elem()
	fn := reflect.MakeFunc(
		reflect.FuncOf([]reflect.Type{inputType}, nil, false),
		func(args []reflect.Value) []reflect.Value {
			input.Set(args[0])
			return nil
		},
	)

	if err := container.Invoke(fn.Interface()); err != nil {
		return *new(T), syserrors.Newf("dig invoke failed for group %q: %w", groupName, err)
	}

	field := input.FieldByName("Funcs")
	if !field.IsValid() {
		return *new(T), syserrors.Newf("field 'Funcs' not found in input struct for group %q", groupName)
	}

	result, ok := field.Interface().(T)
	if !ok {
		return *new(T), syserrors.Newf(
			"type conversion failed for group %q: got %T, want %T",
			groupName,
			field.Interface(),
			*new(T),
		)
	}

	return result, nil
}
