package dig

import (
	"fmt"
	"reflect"
	"strings"

	"frisboo-bank/pkg/syserrors"

	"go.uber.org/dig"
)

// resolveDynamicGroup dynamically resolves a dig group into a value of type T.
// It constructs a struct input with a group tag and invokes the container.
// It returns the resolved group or an error if resolution fails.
func resolveDynamicGroup[T any](container *dig.Container, groupName string, sliceType reflect.Type) (T, error) {
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
		return *new(T), syserrors.Wrapf(err, "dig invoke failed for group %s", groupName)
	}

	field := input.FieldByName("Funcs")
	if !field.IsValid() {
		return *new(T), syserrors.Newf("field 'Funcs' not found in input struct for group %q", groupName)
	}

	result, ok := field.Interface().(T)
	if !ok {
		return *new(T), syserrors.Newf("type conversion failed for group %q: got %T, want %T", groupName, field.Interface(), *new(T))
	}

	return result, nil
}

// wrapFuncAddDigInAndNamed wraps a single-struct-param provider constructor so that:
//  1. The parameter is replaced by an internally generated struct embedding dig.In.
//  2. All exported fields of the original struct become fields of the internal struct.
//  3. If a field has a tag name:"symbolic" and namedDeps maps symbolic -> concrete,
//     the tag is rewritten to name:"concrete".
//  4. If the original parameter already anonymously embeds dig.In (rare, since we enforce plain struct),
//     wrapping is skipped.
//
// The function enforces the provider-level rule: fn must be func(Struct) (T, error).
func wrapFuncWithDigIn(fn any, namedDeps map[string]string, ctx string) (any, error) {
	if fn == nil {
		return nil, syserrors.Newf("%s: function is nil", ctx)
	}

	origVal := reflect.ValueOf(fn)
	origType := origVal.Type()
	if origType.Kind() != reflect.Func {
		return nil, syserrors.Newf("%s: fn not a function", ctx)
	}

	numIn := origType.NumIn()
	if numIn == 0 {
		return fn, nil
	}
	if origType.NumIn() > 1 {
		return nil, syserrors.Newf("%s: constructor must have exactly 1 param (struct), got %d", ctx, origType.NumIn())
	}
	paramType := origType.In(0)
	if paramType.Kind() == reflect.Pointer {
		return nil, syserrors.Newf("%s: param must be a non-pointer struct", ctx)
	}
	if paramType.Kind() != reflect.Struct {
		return nil, syserrors.Newf("%s: param must be struct, got %s", ctx, paramType.Kind())
	}

	// Skip wrapping if already embeds dig.In
	for i := 0; i < paramType.NumField(); i++ {
		sf := paramType.Field(i)
		if sf.Anonymous && sf.Type == reflect.TypeOf(dig.In{}) {
			return fn, nil
		}
	}

	// Build new parameter struct: dig.In + exported fields of original struct.
	internalFields := make([]reflect.StructField, 0, paramType.NumField()+1)
	internalFields = append(internalFields, reflect.StructField{
		Name:      "In",
		Type:      reflect.TypeOf(dig.In{}),
		Anonymous: true,
	})

	type fieldMap struct{ origIndex int }
	mapping := []fieldMap{}

	for i := 0; i < paramType.NumField(); i++ {
		sf := paramType.Field(i)
		if !sf.IsExported() {
			continue
		}
		newTag := string(sf.Tag)
		// If there's a name tag, rewrite if mapping is provided.
		if tagVal, ok := sf.Tag.Lookup("name"); ok {
			actual := namedDeps[tagVal]
			if actual == "" {
				// If no explicit mapping, keep tagVal (so it still uses symbolic).
				actual = tagVal
			}
			// Replace only the name:"..." portion; simplest approach is to rebuild a minimal tag.
			// If you need to preserve other tags, you'd parse them. For now we only care about name.
			newTag = mergeOrReplaceNameTag(string(sf.Tag), actual)
		}

		internalFields = append(internalFields, reflect.StructField{
			Name: sf.Name,
			Type: sf.Type,
			Tag:  reflect.StructTag(newTag),
		})
		mapping = append(mapping, fieldMap{origIndex: i})
	}

	internalParamType := reflect.StructOf(internalFields)

	// Prepare wrapper signature
	numOut := origType.NumOut()
	outTypes := make([]reflect.Type, numOut)
	for i := range numOut {
		outTypes[i] = origType.Out(i)
	}

	wrapperType := reflect.FuncOf([]reflect.Type{internalParamType}, outTypes, false)

	wrapped := reflect.MakeFunc(wrapperType, func(args []reflect.Value) []reflect.Value {
		inVal := args[0]
		// Re-create original param struct value
		origParam := reflect.New(paramType).Elem()
		for idx, m := range mapping {
			fieldVal := inVal.Field(idx + 1) // +1 to skip dig.In
			origParam.Field(m.origIndex).Set(fieldVal)
		}

		results := origVal.Call([]reflect.Value{origParam})
		return results
	})

	return wrapped.Interface(), nil
}

// mergeOrReplaceNameTag attempts to preserve other tag components while replacing or adding name:"value".
// For simplicity, if other tags exist they are preserved as-is and name:"..." is appended or replaced.
func mergeOrReplaceNameTag(existing string, newName string) string {
	if existing == "" {
		return fmt.Sprintf(`name:"%s"`, newName)
	}

	parts := strings.Split(existing, " ")
	found := false
	for i, p := range parts {
		if strings.HasPrefix(p, "name:\"") {
			parts[i] = fmt.Sprintf(`name:"%s"`, newName)
			found = true
			break
		}
	}
	if !found {
		parts = append(parts, fmt.Sprintf(`name:"%s"`, newName))
	}
	return strings.Join(parts, " ")
}
