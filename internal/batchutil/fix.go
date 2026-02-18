package batchutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// KnownTopLevelProps are design properties that belong inside "params", not at op root.
// Exported so validate.go can reference the same list.
var KnownTopLevelProps = []string{
	"x", "y", "width", "height", "w", "h",
	"color", "fillColor", "bg",
	"fontSize", "fontFamily", "fontStyle", "sz", "ff", "fs",
	"parentId", "pid",
	"cornerRadius", "r",
	"layoutMode", "itemSpacing", "padding", "opacity",
	"text", "imageData", "stroke", "strokeWidth",
}

// StripMarkdownFences removes ```json ... ``` or ``` ... ``` wrappers that
// models add even when explicitly told not to.
func StripMarkdownFences(data []byte) []byte {
	s := strings.TrimSpace(string(data))
	if !strings.HasPrefix(s, "```") {
		return data
	}
	lines := strings.Split(s, "\n")
	var filtered []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if i == 0 && strings.HasPrefix(trimmed, "```") {
			continue
		}
		if i == len(lines)-1 && trimmed == "```" {
			continue
		}
		filtered = append(filtered, line)
	}
	return []byte(strings.Join(filtered, "\n"))
}

// FixBatchOps applies auto-corrections to common LLM output drift.
// Returns: fixed JSON bytes, list of human-readable fix descriptions, error.
// Error is non-nil only if the input cannot be parsed at all.
func FixBatchOps(data []byte) ([]byte, []string, error) {
	data = StripMarkdownFences(data)

	// Unwrap {"ops": [...]} or any single-key dict wrapping an array.
	// Models often output {"ops": [...]} instead of bare [...].
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, nil, err
	}
	var fixes []string
	if obj, isObj := raw.(map[string]interface{}); isObj {
		for k, v := range obj {
			if arr, isArr := v.([]interface{}); isArr {
				fixes = append(fixes, fmt.Sprintf("unwrapped dict key %q to get ops array", k))
				b, _ := json.Marshal(arr)
				data = b
				break
			}
		}
	}

	var ops []map[string]interface{}
	if err := json.Unmarshal(data, &ops); err != nil {
		return nil, fixes, err
	}
	for i, op := range ops {
		label := fmt.Sprintf("op[%d]", i)
		if name, ok := op["name"].(string); ok && name != "" {
			label = fmt.Sprintf("op[%d] %q", i, name)
		}

		// Fix "type" → "command"
		if typeVal, hasType := op["type"]; hasType {
			if _, hasCmd := op["command"]; !hasCmd {
				op["command"] = typeVal
				fixes = append(fixes, label+`: renamed "type" to "command"`)
			}
			delete(op, "type")
		}

		// Ensure params exists
		params, _ := op["params"].(map[string]interface{})
		if params == nil {
			params = map[string]interface{}{}
		}

		// Hoist known top-level design props into params
		var hoisted []string
		for _, prop := range KnownTopLevelProps {
			if val, ok := op[prop]; ok {
				params[prop] = val
				delete(op, prop)
				hoisted = append(hoisted, prop)
			}
		}
		if len(hoisted) > 0 {
			fixes = append(fixes, fmt.Sprintf(`%s: moved %s into "params"`, label, strings.Join(hoisted, ", ")))
		}

		op["params"] = params
		ops[i] = op
	}

	fixed, err := json.MarshalIndent(ops, "", "  ")
	if err != nil {
		return nil, fixes, err
	}
	// Preserve ${{...}} interpolation tokens — json.Marshal HTML-escapes < > & by default
	fixed = bytes.ReplaceAll(fixed, []byte(`\u0026`), []byte(`&`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003c`), []byte(`<`))
	fixed = bytes.ReplaceAll(fixed, []byte(`\u003e`), []byte(`>`))
	return fixed, fixes, nil
}
