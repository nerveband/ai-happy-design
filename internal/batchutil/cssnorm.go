package batchutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// NormalizeCSSProps maps CSS-like properties to Figma-native properties.
// It is called during batch parameter normalization so LLMs can write
// familiar CSS properties. CSS props are lower priority — they never
// overwrite Figma-native properties that are already set.
// The function modifies params in-place AND returns it.
func NormalizeCSSProps(params map[string]interface{}) map[string]interface{} {
	if params == nil {
		return params
	}

	// flexDirection → layoutMode
	if fd, ok := params["flexDirection"]; ok {
		if !hasKey(params, "layoutMode") {
			switch fmt.Sprintf("%v", fd) {
			case "column":
				params["layoutMode"] = "VERTICAL"
			case "row":
				params["layoutMode"] = "HORIZONTAL"
			}
		}
		delete(params, "flexDirection")
	}

	// display: "flex" → remove (implied by layoutMode)
	if d, ok := params["display"]; ok {
		if fmt.Sprintf("%v", d) == "flex" {
			delete(params, "display")
		}
	}

	// gap → itemSpacing
	if g, ok := params["gap"]; ok {
		if !hasKey(params, "itemSpacing") {
			params["itemSpacing"] = toFloat(g)
		}
		delete(params, "gap")
	}

	// background / backgroundColor → color
	for _, cssKey := range []string{"background", "backgroundColor"} {
		if bg, ok := params[cssKey]; ok {
			if !hasKey(params, "color") {
				params["color"] = bg
			}
			delete(params, cssKey)
		}
	}

	// borderRadius → cornerRadius
	if br, ok := params["borderRadius"]; ok {
		if !hasKey(params, "cornerRadius") {
			params["cornerRadius"] = toFloat(br)
		}
		delete(params, "borderRadius")
	}

	// padding shorthand
	normalizePadding(params)

	// border shorthand
	normalizeBorder(params)

	// justifyContent → primaryAxisAlignItems
	if jc, ok := params["justifyContent"]; ok {
		if !hasKey(params, "primaryAxisAlignItems") {
			if mapped := mapAlignment(fmt.Sprintf("%v", jc)); mapped != "" {
				params["primaryAxisAlignItems"] = mapped
			}
		}
		delete(params, "justifyContent")
	}

	// alignItems → counterAxisAlignItems
	if ai, ok := params["alignItems"]; ok {
		if !hasKey(params, "counterAxisAlignItems") {
			if mapped := mapAlignment(fmt.Sprintf("%v", ai)); mapped != "" {
				params["counterAxisAlignItems"] = mapped
			}
		}
		delete(params, "alignItems")
	}

	// width: "100%" → layoutSizingHorizontal: "FILL"
	// width: "auto" → layoutSizingHorizontal: "HUG"
	normalizeSizePercent(params, "width", "layoutSizingHorizontal")

	// height: "100%" → layoutSizingVertical: "FILL"
	// height: "auto" → layoutSizingVertical: "HUG"
	normalizeSizePercent(params, "height", "layoutSizingVertical")

	// boxShadow → _shadows
	normalizeBoxShadow(params)

	return params
}

func mapAlignment(css string) string {
	switch css {
	case "flex-start":
		return "MIN"
	case "center":
		return "CENTER"
	case "flex-end":
		return "MAX"
	case "space-between":
		return "SPACE_BETWEEN"
	default:
		return ""
	}
}

func normalizePadding(params map[string]interface{}) {
	p, ok := params["padding"]
	if !ok {
		return
	}
	delete(params, "padding")

	var values []float64

	switch v := p.(type) {
	case float64:
		values = []float64{v}
	case int:
		values = []float64{float64(v)}
	case string:
		parts := strings.Fields(strings.TrimSpace(v))
		for _, part := range parts {
			f, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return // unparseable, bail
			}
			values = append(values, f)
		}
	default:
		// Try numeric conversion
		f := toFloat(p)
		if f != nil {
			values = []float64{f.(float64)}
		}
		return
	}

	if len(values) == 0 {
		return
	}

	var top, right, bottom, left float64
	switch len(values) {
	case 1:
		top, right, bottom, left = values[0], values[0], values[0], values[0]
	case 2:
		top, bottom = values[0], values[0]
		right, left = values[1], values[1]
	case 3:
		top = values[0]
		right, left = values[1], values[1]
		bottom = values[2]
	case 4:
		top, right, bottom, left = values[0], values[1], values[2], values[3]
	default:
		return
	}

	if !hasKey(params, "paddingTop") {
		params["paddingTop"] = top
	}
	if !hasKey(params, "paddingRight") {
		params["paddingRight"] = right
	}
	if !hasKey(params, "paddingBottom") {
		params["paddingBottom"] = bottom
	}
	if !hasKey(params, "paddingLeft") {
		params["paddingLeft"] = left
	}
}

var borderRe = regexp.MustCompile(`^(\d+(?:\.\d+)?)px\s+\w+\s+(.+)$`)

func normalizeBorder(params map[string]interface{}) {
	b, ok := params["border"]
	if !ok {
		return
	}
	delete(params, "border")

	s, ok := b.(string)
	if !ok {
		return
	}

	matches := borderRe.FindStringSubmatch(strings.TrimSpace(s))
	if matches == nil {
		return
	}

	weight, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}

	if !hasKey(params, "strokeWeight") {
		params["strokeWeight"] = weight
	}
	if !hasKey(params, "stroke") {
		params["stroke"] = strings.TrimSpace(matches[2])
	}
}

func normalizeSizePercent(params map[string]interface{}, cssKey, figmaKey string) {
	v, ok := params[cssKey]
	if !ok {
		return
	}

	s, ok := v.(string)
	if !ok {
		return
	}

	if hasKey(params, figmaKey) {
		return
	}

	switch s {
	case "100%":
		params[figmaKey] = "FILL"
		delete(params, cssKey)
	case "auto":
		params[figmaKey] = "HUG"
		delete(params, cssKey)
	}
}

// shadowRe captures: offsetX offsetY blur [spread] color
// Values may have optional "px" suffix. Color can be rgba(...) or hex with optional alpha.
var shadowRe = regexp.MustCompile(
	`(-?\d+(?:\.\d+)?)(?:px)?\s+(-?\d+(?:\.\d+)?)(?:px)?\s+(\d+(?:\.\d+)?)(?:px)?\s+` +
		`(?:(\d+(?:\.\d+)?)(?:px)?\s+)?` +
		`(rgba?\([^)]+\)|#[0-9a-fA-F]+)`,
)

func normalizeBoxShadow(params map[string]interface{}) {
	bs, ok := params["boxShadow"]
	if !ok {
		return
	}
	delete(params, "boxShadow")

	s, ok := bs.(string)
	if !ok {
		return
	}

	if hasKey(params, "_shadows") {
		return
	}

	// Split on commas that are NOT inside parentheses.
	shadows := splitShadows(s)
	var result []map[string]interface{}

	for _, shadow := range shadows {
		shadow = strings.TrimSpace(shadow)
		matches := shadowRe.FindStringSubmatch(shadow)
		if matches == nil {
			continue
		}

		ox, _ := strconv.ParseFloat(matches[1], 64)
		oy, _ := strconv.ParseFloat(matches[2], 64)
		r, _ := strconv.ParseFloat(matches[3], 64)
		color := matches[5]

		entry := map[string]interface{}{
			"offsetX": ox,
			"offsetY": oy,
			"radius":  r,
			"color":   color,
		}

		if matches[4] != "" {
			spread, _ := strconv.ParseFloat(matches[4], 64)
			entry["spread"] = spread
		}

		result = append(result, entry)
	}

	if len(result) > 0 {
		params["_shadows"] = result
	}
}

// splitShadows splits on commas outside parentheses.
func splitShadows(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, ch := range s {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
}

// toFloat attempts to convert a value to float64.
// Returns float64 if successful, or the original value if not.
func toFloat(v interface{}) interface{} {
	switch val := v.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
		if err != nil {
			return v
		}
		return f
	default:
		return v
	}
}
