package enum

import (
	"fmt"
	"reflect"

	"github.com/stoewer/go-strcase"
)

type definition[T comparable] struct {
	values map[T]string
	names  map[string]T
	zero   T
}

func newDefinition[T comparable](mapping map[T]string) *definition[T] {
	d := &definition[T]{
		values: mapping,
		names:  make(map[string]T, len(mapping)),
	}

	var zero T
	d.zero = zero

	for value, name := range mapping {
		name = strcase.UpperCamelCase(name)
		d.names[name] = value

		if value == zero {
			d.zero = value
		}
	}

	return d
}

// String returns the string name for a value, or UNKNOWN(n) if not found.
func (d *definition[T]) String(val T) string {
	if name, ok := d.values[val]; ok {
		return name
	}

	return fmt.Sprintf("UNKNOWN(%v)", val)
}

// Valid returns true if the value is defined in the enum.
func (d *definition[T]) Valid(val T) bool {
	_, ok := d.values[val]
	return ok
}

// Parse converts input to enum value
func (d *definition[T]) Parse(val any) (T, error) {
	switch val := val.(type) {
	case string:
		val = strcase.UpperCamelCase(val)

		if match, ok := d.names[val]; ok {
			return match, nil
		}
	case T:
		if d.Valid(val) {
			return val, nil
		}
	}

	// Try conversion for int types if T is int-like
	switch any(d.zero).(type) {
	case int, int8, int16, int32, int64:
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if rv.Type().ConvertibleTo(reflect.TypeOf(d.zero)) {
				converted := rv.Convert(reflect.TypeOf(d.zero)).Interface().(T)
				if d.Valid(converted) {
					return converted, nil
				}
			}
		}
	}

	return d.zero, fmt.Errorf("invalid enum value: %v", val)
}

// MustParse parses input and panics on error
func (d *definition[T]) MustParse(val any) T {
	match, err := d.Parse(val)
	if err != nil {
		panic(err)
	}

	return match
}

// IsZero checks if value is the zero value
func (d *definition[T]) IsZero(val T) bool {
	return val == d.zero
}

// Names returns all enum names.
func (d *definition[T]) Names() []string {
	out := make([]string, 0, len(d.names))

	for name := range d.names {
		out = append(out, name)
	}

	return out
}

// Name returns the name for a value, or false if missing.
func (d *definition[T]) Name(val T) (string, bool) {
	name, ok := d.values[val]
	return name, ok
}

// FromName looks up a value by name (case-insensitive, normalized).
func (d *definition[T]) FromName(name string) (T, error) {
	key := strcase.UpperCamelCase(name)
	if val, ok := d.names[key]; ok {
		return val, nil
	}

	return d.zero, fmt.Errorf("invalid enum name: %s", name)
}

// Values returns all defined values
func (d *definition[T]) Values() []T {
	result := make([]T, 0, len(d.values))

	for value := range d.values {
		result = append(result, value)
	}

	return result
}
