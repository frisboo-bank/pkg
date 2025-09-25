package dig

import (
	"fmt"
	"reflect"

	"frisboo-bank/pkg/syserrors"

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
