package figmacli

import (
	"encoding/json"
	"fmt"
	"strings"
)

func JSONEnvelope(command string, result interface{}, meta map[string]interface{}) ResultEnvelope {
	if meta == nil {
		meta = map[string]interface{}{}
	}
	return ResultEnvelope{OK: true, Command: command, Result: result, Meta: meta}
}

func ErrorJSON(command, code, message, hint string, retryable bool) ResultEnvelope {
	return ResultEnvelope{
		OK:      false,
		Command: command,
		Error:   &ErrorEnvelope{Code: code, Message: message, Hint: hint, Retryable: retryable},
	}
}

func Marshal(format string, v interface{}) ([]byte, error) {
	switch format {
	case "", "json":
		return json.MarshalIndent(v, "", "  ")
	case "jsonl":
		data, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		return append(data, '\n'), nil
	case "text":
		if s, ok := v.(string); ok {
			return []byte(s + "\n"), nil
		}
		return json.MarshalIndent(v, "", "  ")
	default:
		return nil, fmt.Errorf("unsupported output format %q", format)
	}
}

func ParseFields(fields string) []string {
	parts := strings.Split(fields, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}
