package commoncli

import "strings"

// FilterFields returns a shallow filtered view of a JSON-like object.
func FilterFields(value any, fields string) any {
	trimmed := strings.TrimSpace(fields)
	if trimmed == "" {
		return value
	}
	obj, ok := value.(map[string]any)
	if !ok {
		return value
	}
	out := map[string]any{}
	for _, field := range strings.Split(trimmed, ",") {
		key := strings.TrimSpace(field)
		if key == "" {
			continue
		}
		if v, exists := obj[key]; exists {
			out[key] = v
		}
	}
	return out
}
