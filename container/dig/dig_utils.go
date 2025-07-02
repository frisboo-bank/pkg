package dig

import (
	"fmt"
	"reflect"

	loggerContracts "frisboo-bank/pkg/logger/contracts"

	"go.uber.org/dig"
)

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
		return *new(T), err
	}

	field := input.FieldByName("Funcs")
	if !field.IsValid() {
		return *new(T), fmt.Errorf("field 'Funcs' not found in input struct")
	}

	result, ok := field.Interface().(T)
	if !ok {
		return *new(T), fmt.Errorf("type conversion failed for group %s: got %T, want %T", groupName, field, *new(T))
	}

	return result, nil
}

func filterOptions[T any, O any](options []O, logger loggerContracts.Logger) []T {
	var filteredOptions []T

	for _, option := range options {
		if opt, ok := any(option).(T); !ok {
			logger.Errorw("(dig-container) option doesn't seem to be a valid dig option",
				loggerContracts.Fields{
					"option":   option,
					"expected": fmt.Sprintf("%T", *new(T)),
				})
		} else {
			filteredOptions = append(filteredOptions, opt)
		}
	}

	return filteredOptions
}
