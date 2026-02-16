package batchutil

import (
	"regexp"
	"strings"
)

// NormalizeBatchParams expands shorthand keys for batch operations.
// It is intentionally command-aware to avoid key collisions.
func NormalizeBatchParams(command string, params map[string]interface{}) map[string]interface{} {
	out := cloneParams(params)
	if len(out) == 0 {
		return out
	}

	applyAlias(out, "width", "w")
	applyAlias(out, "height", "h")
	applyAlias(out, "parentId", "pid")

	// Universal fillColor → color alias (the #1 silent failure from cross-tool training data)
	applyAlias(out, "color", "fillColor")

	cmd := normalizeCommand(command)

	if isTextCommand(cmd) {
		applyAlias(out, "fontSize", "sz")
		applyAlias(out, "fontFamily", "ff")
		applyAlias(out, "fontStyle", "fs")
		applyAlias(out, "lineHeight", "lh")
		applyAlias(out, "letterSpacing", "ls")
		if hasKey(out, "lineHeight") && !hasKey(out, "lineHeightUnit") {
			out["lineHeightUnit"] = "PERCENT"
		}
	}

	if isFrameOrRectCreateCommand(cmd) || isModifyCommand(cmd) {
		applyAlias(out, "color", "bg")
		applyAlias(out, "cornerRadius", "r")
	}

	if isShapeOrFrameCreateCommand(cmd) {
		applyAlias(out, "strokeWidth", "sw")
	}

	if isPaintStrokeCommand(cmd) {
		applyAlias(out, "strokeWeight", "sw")
	}

	return out
}

// stepNameRe matches characters that are NOT alphanumeric or underscore.
var stepNameRe = regexp.MustCompile(`[^a-z0-9_]+`)

// SanitizeStepName normalizes a batch step name: lowercase, spaces→underscores,
// strip non-alphanumeric characters (except _).
func SanitizeStepName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = stepNameRe.ReplaceAllString(name, "")
	// Collapse multiple underscores
	for strings.Contains(name, "__") {
		name = strings.ReplaceAll(name, "__", "_")
	}
	name = strings.Trim(name, "_")
	return name
}

func normalizeCommand(command string) string {
	return strings.ToLower(strings.TrimSpace(command))
}

func cloneParams(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		return map[string]interface{}{}
	}
	out := make(map[string]interface{}, len(params))
	for k, v := range params {
		out[k] = v
	}
	return out
}

func hasKey(params map[string]interface{}, key string) bool {
	_, ok := params[key]
	return ok
}

func applyAlias(params map[string]interface{}, canonical, shorthand string) {
	if hasKey(params, canonical) {
		return
	}
	if val, ok := params[shorthand]; ok {
		params[canonical] = val
	}
}

func isTextCommand(cmd string) bool {
	switch cmd {
	case "text", "text.create", "text.set_spacing", "text.set_line_height", "text.set_letter_spacing",
		"create_text", "set_text_spacing", "set_line_height", "set_letter_spacing", "set_paragraph_spacing":
		return true
	default:
		return false
	}
}

func isFrameOrRectCreateCommand(cmd string) bool {
	switch cmd {
	case "frame", "rect", "image",
		"node.create_frame", "shape.create_rectangle", "shape.create_image",
		"create_frame", "create_rectangle", "create_image":
		return true
	default:
		return false
	}
}

func isShapeOrFrameCreateCommand(cmd string) bool {
	switch cmd {
	case "frame", "rect", "ellipse", "line", "image",
		"node.create_frame", "shape.create_rectangle", "shape.create_ellipse", "shape.create_polygon", "shape.create_star", "shape.create_line", "shape.create_image",
		"create_frame", "create_rectangle", "create_ellipse", "create_polygon", "create_star", "create_line", "create_image":
		return true
	default:
		return false
	}
}

func isPaintStrokeCommand(cmd string) bool {
	switch cmd {
	case "stroke", "paint.set_stroke", "set_stroke", "set_stroke_color":
		return true
	default:
		return false
	}
}

func isModifyCommand(cmd string) bool {
	switch cmd {
	case "modify", "node.modify":
		return true
	default:
		return false
	}
}
