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

func wrapFuncWithNamedInputs(fn any, namedDeps map[string]string, ctx string) (any, error) {
	if fn == nil {
		return nil, syserrors.Newf("%s: function is nil", ctx)
	}
	if len(namedDeps) == 0 {
		return fn, nil
	}

	origVal := reflect.ValueOf(fn)
	origType := origVal.Type()
	if origType.Kind() != reflect.Func {
		return nil, syserrors.Newf("%s: fn not a function", ctx)
	}

	numIn := origType.NumIn()
	if numIn != 1 {
		return fn, nil
	}

	paramType := origType.In(0)
	ptrParam := false
	if paramType.Kind() == reflect.Pointer {
		paramType = paramType.Elem()
		ptrParam = true
	}
	if paramType.Kind() != reflect.Struct {
		return fn, nil
	}

	// Collect exported struct fields with name tag.
	var adaptedFields []reflect.StructField
	var originalFieldIndexes []int
	for i := range paramType.NumField() {
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
		adaptedFields = append(adaptedFields, reflect.StructField{
			Name: sf.Name,
			Type: sf.Type,
			Tag:  reflect.StructTag(`name:"` + actual + `"`),
		})
		originalFieldIndexes = append(originalFieldIndexes, i)
	}

	// If no fields need adaptation, use original.
	if len(adaptedFields) == 0 {
		return fn, nil
	}

	// Build internal dig.In struct type.
	fields := []reflect.StructField{
		{
			Name:      "In",
			Type:      reflect.TypeOf(dig.In{}),
			Anonymous: true,
		},
	}
	fields = append(fields, adaptedFields...)
	internalParamType := reflect.StructOf(fields)

	// Outputs unchanged.
	numOut := origType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for i := range numOut {
		outTypes[i] = origType.Out(i)
	}

	wrapperType := reflect.FuncOf([]reflect.Type{internalParamType}, outTypes, false)
	wrapped := reflect.MakeFunc(wrapperType, func(args []reflect.Value) []reflect.Value {
		inVal := args[0]

		var origParam reflect.Value
		if ptrParam {
			origParam = reflect.New(paramType)
		} else {
			origParam = reflect.New(paramType).Elem()
		}

		for idx, origFieldIdx := range originalFieldIndexes {
			fieldVal := inVal.Field(idx + 1) // skip dig.In at index 0
			if ptrParam {
				origParam.Elem().Field(origFieldIdx).Set(fieldVal)
			} else {
				origParam.Field(origFieldIdx).Set(fieldVal)
			}
		}

		callArgs := []reflect.Value{origParam}
		return origVal.Call(callArgs)
	})

	return wrapped.Interface(), nil
}
