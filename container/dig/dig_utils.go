package dig

import (
	"fmt"
	"reflect"

	"go.uber.org/dig"
)

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
		return *new(T), fmt.Errorf("type conversion failed for group %q: got %T, want %T", groupName, field.Interface(), *new(T))
	}

	return result, nil
}

// filterOptions filters a slice of options of type O to a slice of type T.
// Returns an error if any option cannot be converted to T.
func filterOptions[T any, O any](options []O) ([]T, error) {
	var filteredOptions []T

	for idx, option := range options {
		opt, ok := any(option).(T)
		if !ok {
			return nil, fmt.Errorf("option at index %d must be of type %T but is currently of type %T", idx, *new(T), option)
		}
		filteredOptions = append(filteredOptions, opt)
	}

	return filteredOptions, nil
}
