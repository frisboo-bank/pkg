package utils

import "fmt"

// Helper: prepend prefix string to message (for ...any signatures)
func ConcatPrefix(prefix string, v ...any) []any {
	if prefix == "" {
		return v
	}

	return append([]any{prefix}, v...)
}

// Helper: prepend prefix to format string (for *f methods)
func ConcatPrefixf(prefix string, format string, v ...any) (string, []any) {
	if prefix == "" {
		return format, v
	}

	return fmt.Sprintf("%s %s", prefix, format), v
}

// Helper: prepend prefix to string message (for *w methods)
func ConcatPrefixStr(prefix string, message string) string {
	if prefix == "" {
		return message
	}
	return fmt.Sprintf("%s %s", prefix, message)
}
