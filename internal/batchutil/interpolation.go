package batchutil

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`\$\{\{\s*([^{}]+?)\s*\}\}`)

// StepState captures the execution state of a prior batch step.
type StepState struct {
	Index   int
	Name    string
	Command string
	OK      bool
	Result  interface{}
	Error   string
}

// BuildContext creates an interpolation context from completed steps.
func BuildContext(steps []StepState) map[string]interface{} {
	stepsMap := make(map[string]interface{}, len(steps)*2)
	var last map[string]interface{}

	for _, step := range steps {
		entry := map[string]interface{}{
			"index":   step.Index,
			"name":    step.Name,
			"command": step.Command,
			"ok":      step.OK,
			"result":  step.Result,
			"error":   step.Error,
		}
		stepsMap[strconv.Itoa(step.Index)] = entry
		if step.Name != "" {
			stepsMap[step.Name] = entry
		}
		last = entry
	}

	ctx := map[string]interface{}{
		"steps": stepsMap,
	}
	if last != nil {
		ctx["last"] = last
	}
	return ctx
}

// InterpolateParams resolves placeholders in params using prior step results.
func InterpolateParams(params map[string]interface{}, steps []StepState) (map[string]interface{}, error) {
	ctx := BuildContext(steps)
	value, err := InterpolateValue(params, ctx)
	if err != nil {
		return nil, err
	}
	out, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("interpolation produced invalid params type %T", value)
	}
	return out, nil
}

// InterpolateValue recursively resolves placeholders in arbitrary JSON-like data.
func InterpolateValue(value interface{}, ctx map[string]interface{}) (interface{}, error) {
	switch typed := value.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{}, len(typed))
		for key, item := range typed {
			resolved, err := InterpolateValue(item, ctx)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
			out[key] = resolved
		}
		return out, nil
	case []interface{}:
		out := make([]interface{}, len(typed))
		for i, item := range typed {
			resolved, err := InterpolateValue(item, ctx)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %w", i, err)
			}
			out[i] = resolved
		}
		return out, nil
	case string:
		return interpolateString(typed, ctx)
	default:
		return value, nil
	}
}

func interpolateString(input string, ctx map[string]interface{}) (interface{}, error) {
	input = expandShortInterpolation(input)

	matches := placeholderPattern.FindAllStringSubmatchIndex(input, -1)
	if len(matches) == 0 {
		return input, nil
	}

	if len(matches) == 1 && matches[0][0] == 0 && matches[0][1] == len(input) {
		expr := strings.TrimSpace(input[matches[0][2]:matches[0][3]])
		value, err := resolvePath(ctx, expr)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	var builder strings.Builder
	last := 0
	for _, match := range matches {
		builder.WriteString(input[last:match[0]])
		expr := strings.TrimSpace(input[match[2]:match[3]])
		value, err := resolvePath(ctx, expr)
		if err != nil {
			return nil, err
		}
		builder.WriteString(stringify(value))
		last = match[1]
	}
	builder.WriteString(input[last:])
	return builder.String(), nil
}

func expandShortInterpolation(input string) string {
	if input == "" {
		return input
	}

	var b strings.Builder
	b.Grow(len(input) + 16)

	for i := 0; i < len(input); {
		ch := input[i]
		if ch != '$' {
			b.WriteByte(ch)
			i++
			continue
		}

		// Keep long-form interpolation untouched.
		if i+1 < len(input) && input[i+1] == '{' {
			b.WriteByte(ch)
			i++
			continue
		}

		// "$name" shorthand should not trigger in embedded identifiers like "foo$bar".
		if i > 0 && isIdentChar(input[i-1]) {
			b.WriteByte(ch)
			i++
			continue
		}

		start := i + 1
		if start >= len(input) || !isIdentStart(input[start]) {
			b.WriteByte(ch)
			i++
			continue
		}

		j := start + 1
		for j < len(input) && isIdentChar(input[j]) {
			j++
		}

		base := input[start:j]
		path := ""
		for j < len(input) && input[j] == '.' {
			segStart := j + 1
			if segStart >= len(input) || !isIdentStart(input[segStart]) {
				break
			}
			j = segStart + 1
			for j < len(input) && isIdentChar(input[j]) {
				j++
			}
			if path == "" {
				path = input[segStart:j]
			} else {
				path += "." + input[segStart:j]
			}
		}

		if strings.EqualFold(base, "last") {
			if path == "" {
				b.WriteString("${{last.result.id}}")
			} else {
				b.WriteString("${{last.result.")
				b.WriteString(path)
				b.WriteString("}}")
			}
		} else {
			b.WriteString("${{steps.")
			b.WriteString(base)
			b.WriteString(".result.")
			if path == "" {
				b.WriteString("id")
			} else {
				b.WriteString(path)
			}
			b.WriteString("}}")
		}
		i = j
	}

	return b.String()
}

func isIdentStart(ch byte) bool {
	return ch == '_' || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
}

func isIdentChar(ch byte) bool {
	return isIdentStart(ch) || (ch >= '0' && ch <= '9')
}

func resolvePath(root interface{}, path string) (interface{}, error) {
	normalized := normalizePath(path)
	if normalized == "" {
		return nil, fmt.Errorf("empty interpolation path")
	}

	parts := strings.Split(normalized, ".")
	current := root
	for _, part := range parts {
		if part == "" {
			continue
		}
		switch typed := current.(type) {
		case map[string]interface{}:
			next, ok := typed[part]
			if !ok {
				return nil, fmt.Errorf("interpolation path not found: %s", path)
			}
			current = next
		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid list index %q in path %s", part, path)
			}
			if index < 0 || index >= len(typed) {
				return nil, fmt.Errorf("list index out of range %d in path %s", index, path)
			}
			current = typed[index]
		default:
			return nil, fmt.Errorf("cannot resolve %q in path %s", part, path)
		}
	}
	return current, nil
}

func normalizePath(path string) string {
	normalized := strings.TrimSpace(path)
	normalized = strings.ReplaceAll(normalized, "[", ".")
	normalized = strings.ReplaceAll(normalized, "]", "")
	for strings.Contains(normalized, "..") {
		normalized = strings.ReplaceAll(normalized, "..", ".")
	}
	normalized = strings.Trim(normalized, ".")
	return normalized
}

func stringify(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return "null"
	case string:
		return typed
	case float64, float32, int, int32, int64, uint, uint32, uint64, bool:
		return fmt.Sprintf("%v", typed)
	default:
		data, err := json.Marshal(typed)
		if err != nil {
			return fmt.Sprintf("%v", typed)
		}
		return string(data)
	}
}
