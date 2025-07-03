package enum

// Enum is a generic enum wrapper for a given set of values.
type Enum[T comparable] struct {
	*definition[T]
}

// New creates a new Enum for a given mapping.
func New[T comparable](mapping map[T]string) *Enum[T] {
	return &Enum[T]{
		definition: newDefinition[T](mapping),
	}
}

func (e *Enum[T]) String(val T) string             { return e.definition.String(val) }
func (e *Enum[T]) Valid(val T) bool                { return e.definition.Valid(val) }
func (e *Enum[T]) Name(val T) (string, bool)       { return e.definition.Name(val) }
func (e *Enum[T]) Names() []string                 { return e.definition.Names() }
func (e *Enum[T]) FromName(name string) (T, error) { return e.definition.FromName(name) }
func (e *Enum[T]) IsZero(val T) bool               { return e.definition.IsZero(val) }
func (e *Enum[T]) Values() []T                     { return e.definition.Values() }
func (e *Enum[T]) Parse(val any) (T, error)        { return e.definition.Parse(val) }
func (e *Enum[T]) MustParse(val any) T             { return e.definition.MustParse(val) }

// MarshalText implements encoding.TextMarshaler
func (e *Enum[T]) MarshalText(val T) ([]byte, error) {
	return []byte(e.String(val)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *Enum[T]) UnmarshalText(val *T, text []byte) error {
	parsed, err := e.Parse(string(text))
	if err != nil {
		return err
	}

	*val = parsed
	return nil
}
