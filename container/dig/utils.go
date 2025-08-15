package dig

import (
	"fmt"
	"reflect"

	"frisboo-bank/pkg/container/dependencies/decorator"
	"frisboo-bank/pkg/container/dependencies/hook"
	"frisboo-bank/pkg/container/dependencies/invoker"
	"frisboo-bank/pkg/container/dependencies/provider"

	"go.uber.org/dig"
)

func ToDigDecoratorOptions(opts decorator.DecoratorOptions) []dig.DecorateOption {
	var result []dig.DecorateOption

	return result
}

func ToDigHookOptions(opts hook.HooksOptions) []dig.ProvideOption {
	var result []dig.ProvideOption

	if len(opts.As) > 0 {
		result = append(result, dig.As(opts.As...))
	}

	if opts.Export {
		result = append(result, dig.Export(true))
	}

	if opts.Group != "" {
		result = append(result, dig.Group(opts.Group))
	}

	if opts.LocationPC != 0 {
		result = append(result, dig.LocationForPC(opts.LocationPC))
	}

	if opts.Name != "" {
		result = append(result, dig.Name(opts.Name))
	}

	return result
}

func ToDigInvokerOptions(opts invoker.InvokerOptions) []dig.InvokeOption {
	var result []dig.InvokeOption

	return result
}

func ToDigProvideOptions(opts provider.ProviderOptions) []dig.ProvideOption {
	var result []dig.ProvideOption

	if len(opts.As) > 0 {
		result = append(result, dig.As(opts.As...))
	}

	if opts.Export {
		result = append(result, dig.Export(true))
	}

	if opts.Group != "" {
		result = append(result, dig.Group(opts.Group))
	}

	if opts.LocationPC != 0 {
		result = append(result, dig.LocationForPC(opts.LocationPC))
	}

	if opts.Name != "" {
		result = append(result, dig.Name(opts.Name))
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
		return *new(T), fmt.Errorf("dig invoke failed for group %q: %w", groupName, err)
	}

	field := input.FieldByName("Funcs")
	if !field.IsValid() {
		return *new(T), fmt.Errorf("field 'Funcs' not found in input struct for group %q", groupName)
	}

	result, ok := field.Interface().(T)
	if !ok {
		return *new(T), fmt.Errorf(
			"type conversion failed for group %q: got %T, want %T",
			groupName,
			field.Interface(),
			*new(T),
		)
	}

	return result, nil
}
