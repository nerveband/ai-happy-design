package extract

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// GradientStop represents a color stop in a CSS linear-gradient.
type GradientStop struct {
	Color    string  // hex color, e.g. "#0C1E2C"
	Position float64 // 0.0–1.0
}

// parseStyleBlocks extracts CSS class rules from <style> blocks in HTML.
// Returns map[className]map[property]value, e.g. {".bg-dark": {"background": "linear-gradient(...)"}}
func parseStyleBlocks(html string) map[string]map[string]string {
	result := make(map[string]map[string]string)

	// Extract all <style>...</style> blocks
	styleRe := regexp.MustCompile(`(?is)<style[^>]*>(.*?)</style>`)
	for _, match := range styleRe.FindAllStringSubmatch(html, -1) {
		css := match[1]
		parseRules(css, result)
	}

	return result
}

// parseRules parses CSS rules from a CSS string into the result map.
func parseRules(css string, result map[string]map[string]string) {
	// Remove CSS comments
	commentRe := regexp.MustCompile(`/\*.*?\*/`)
	css = commentRe.ReplaceAllString(css, "")

	// Remove @keyframes blocks (they contain braces that confuse the rule parser)
	keyframesRe := regexp.MustCompile(`(?s)@keyframes\s+\w+\s*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}`)
	css = keyframesRe.ReplaceAllString(css, "")

	// Remove @media and other at-rules for simplicity
	atRuleRe := regexp.MustCompile(`(?s)@\w+[^{]*\{[^{}]*(?:\{[^{}]*\}[^{}]*)*\}`)
	css = atRuleRe.ReplaceAllString(css, "")

	// Match selector { declarations }
	// We handle simple selectors: .class, .class.class, element.class, etc.
	ruleRe := regexp.MustCompile(`(?s)([^{}]+?)\{([^{}]*)\}`)
	for _, match := range ruleRe.FindAllStringSubmatch(css, -1) {
		selectors := strings.TrimSpace(match[1])
		declarations := match[2]

		// Parse declarations
		props := parseDeclarations(declarations)
		if len(props) == 0 {
			continue
		}

		// Handle comma-separated selectors
		for _, sel := range strings.Split(selectors, ",") {
			sel = strings.TrimSpace(sel)
			if sel == "" {
				continue
			}
			// Store under the full selector
			if result[sel] == nil {
				result[sel] = make(map[string]string)
			}
			for k, v := range props {
				result[sel][k] = v
			}
		}
	}
}

// parseDeclarations parses "prop: value; prop: value;" into a map.
func parseDeclarations(block string) map[string]string {
	result := make(map[string]string)
	for _, decl := range strings.Split(block, ";") {
		decl = strings.TrimSpace(decl)
		if decl == "" {
			continue
		}
		idx := strings.Index(decl, ":")
		if idx < 0 {
			continue
		}
		prop := strings.TrimSpace(decl[:idx])
		value := strings.TrimSpace(decl[idx+1:])
		if prop != "" && value != "" {
			result[prop] = value
		}
	}
	return result
}

// parseInlineStyle parses a style="..." attribute value.
func parseInlineStyle(style string) map[string]string {
	return parseDeclarations(style)
}

// parseLinearGradient parses CSS linear-gradient() into angle and stops.
// e.g. "linear-gradient(150deg, #0C1E2C 0%, #14344A 100%)"
func parseLinearGradient(value string) (float64, []GradientStop, error) {
	value = strings.TrimSpace(value)

	// Extract the contents of linear-gradient(...)
	re := regexp.MustCompile(`(?i)linear-gradient\s*\((.+)\)`)
	m := re.FindStringSubmatch(value)
	if m == nil {
		return 0, nil, fmt.Errorf("not a linear-gradient: %q", value)
	}
	inner := strings.TrimSpace(m[1])

	// Split by commas, but be careful with nested parens (e.g. rgba())
	parts := splitGradientParts(inner)
	if len(parts) < 2 {
		return 0, nil, fmt.Errorf("too few gradient parts in %q", value)
	}

	var angle float64
	stopStart := 0

	// First part might be an angle or direction
	first := strings.TrimSpace(parts[0])
	if a, ok := parseAngle(first); ok {
		angle = a
		stopStart = 1
	} else if dir, ok := parseDirection(first); ok {
		angle = dir
		stopStart = 1
	} else {
		// No angle, default is 180deg (top to bottom)
		angle = 180
	}

	// Parse color stops
	var stops []GradientStop
	for i := stopStart; i < len(parts); i++ {
		part := strings.TrimSpace(parts[i])
		if part == "" {
			continue
		}
		stop, err := parseColorStop(part)
		if err != nil {
			continue // skip unparseable stops
		}
		stops = append(stops, stop)
	}

	if len(stops) < 2 {
		return 0, nil, fmt.Errorf("need at least 2 color stops, got %d", len(stops))
	}

	// If stops don't have explicit positions, distribute evenly
	distributeStopPositions(stops)

	return angle, stops, nil
}

// splitGradientParts splits gradient contents by commas, respecting nested parentheses.
func splitGradientParts(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, c := range s {
		switch c {
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

// parseAngle parses "150deg" or "1.5rad" or "0.5turn" into degrees.
func parseAngle(s string) (float64, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	if strings.HasSuffix(s, "deg") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "deg"), 64)
		if err == nil {
			return v, true
		}
	}
	if strings.HasSuffix(s, "rad") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "rad"), 64)
		if err == nil {
			return v * 180 / math.Pi, true
		}
	}
	if strings.HasSuffix(s, "turn") {
		v, err := strconv.ParseFloat(strings.TrimSuffix(s, "turn"), 64)
		if err == nil {
			return v * 360, true
		}
	}
	return 0, false
}

// parseDirection parses "to right", "to bottom left", etc. into degrees.
func parseDirection(s string) (float64, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	if !strings.HasPrefix(s, "to ") {
		return 0, false
	}
	dir := strings.TrimPrefix(s, "to ")
	switch dir {
	case "top":
		return 0, true
	case "right":
		return 90, true
	case "bottom":
		return 180, true
	case "left":
		return 270, true
	case "top right":
		return 45, true
	case "bottom right":
		return 135, true
	case "bottom left":
		return 225, true
	case "top left":
		return 315, true
	}
	return 0, false
}

// parseColorStop parses "#0C1E2C 0%" or "rgba(255,0,0,0.5) 50%" into a GradientStop.
func parseColorStop(part string) (GradientStop, error) {
	part = strings.TrimSpace(part)
	stop := GradientStop{Position: -1} // -1 means "not set"

	// Try to find a percentage at the end
	pctRe := regexp.MustCompile(`\s+([\d.]+)%\s*$`)
	if m := pctRe.FindStringSubmatch(part); m != nil {
		pct, err := strconv.ParseFloat(m[1], 64)
		if err == nil {
			stop.Position = pct / 100.0
		}
		part = strings.TrimSpace(part[:len(part)-len(m[0])])
	}

	// The rest is the color
	hex := cssColorToHex(part)
	if hex == "" {
		return stop, fmt.Errorf("cannot parse color %q", part)
	}
	stop.Color = hex
	return stop, nil
}

// distributeStopPositions fills in -1 positions evenly between set positions.
func distributeStopPositions(stops []GradientStop) {
	if len(stops) == 0 {
		return
	}
	// If first has no position, set to 0
	if stops[0].Position < 0 {
		stops[0].Position = 0
	}
	// If last has no position, set to 1
	if stops[len(stops)-1].Position < 0 {
		stops[len(stops)-1].Position = 1
	}
	// Fill in gaps
	for i := 1; i < len(stops)-1; i++ {
		if stops[i].Position < 0 {
			// Find next set position
			next := i + 1
			for next < len(stops) && stops[next].Position < 0 {
				next++
			}
			// Distribute evenly
			prevPos := stops[i-1].Position
			nextPos := stops[next].Position
			count := next - i + 1
			step := (nextPos - prevPos) / float64(count)
			for j := i; j < next; j++ {
				stops[j].Position = prevPos + step*float64(j-i+1)
			}
		}
	}
}

// cssColorToHex converts various CSS color formats to hex (#RRGGBB).
// Handles: hex (#RGB, #RRGGBB, #RRGGBBAA), rgb(), rgba(), named colors.
func cssColorToHex(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	// Hex colors
	if strings.HasPrefix(value, "#") {
		hex := strings.TrimPrefix(value, "#")
		switch len(hex) {
		case 3:
			// #RGB -> #RRGGBB
			return "#" + strings.ToUpper(string(hex[0])+string(hex[0])+string(hex[1])+string(hex[1])+string(hex[2])+string(hex[2]))
		case 4:
			// #RGBA -> #RRGGBB (drop alpha)
			return "#" + strings.ToUpper(string(hex[0])+string(hex[0])+string(hex[1])+string(hex[1])+string(hex[2])+string(hex[2]))
		case 6:
			return "#" + strings.ToUpper(hex)
		case 8:
			// #RRGGBBAA -> #RRGGBB (drop alpha)
			return "#" + strings.ToUpper(hex[:6])
		}
		return "#" + strings.ToUpper(hex)
	}

	// rgb() and rgba()
	rgbRe := regexp.MustCompile(`(?i)rgba?\s*\(\s*([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)`)
	if m := rgbRe.FindStringSubmatch(value); m != nil {
		r := clampInt(parseFloatDef(m[1], 0), 0, 255)
		g := clampInt(parseFloatDef(m[2], 0), 0, 255)
		b := clampInt(parseFloatDef(m[3], 0), 0, 255)
		return fmt.Sprintf("#%02X%02X%02X", r, g, b)
	}

	// Named colors (common subset)
	named := map[string]string{
		"black":       "#000000",
		"white":       "#FFFFFF",
		"red":         "#FF0000",
		"green":       "#008000",
		"blue":        "#0000FF",
		"yellow":      "#FFFF00",
		"cyan":        "#00FFFF",
		"magenta":     "#FF00FF",
		"orange":      "#FFA500",
		"purple":      "#800080",
		"pink":        "#FFC0CB",
		"gray":        "#808080",
		"grey":        "#808080",
		"silver":      "#C0C0C0",
		"navy":        "#000080",
		"teal":        "#008080",
		"transparent": "",
	}
	if hex, ok := named[strings.ToLower(value)]; ok {
		return hex
	}

	return ""
}

// remToPx converts a CSS rem value string to pixels.
// e.g. "1.6rem" with baseFontSize=16 -> 25.6
func remToPx(value string, baseFontSize float64) float64 {
	value = strings.TrimSpace(strings.ToLower(value))
	if strings.HasSuffix(value, "rem") {
		num := strings.TrimSuffix(value, "rem")
		v, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
		if err == nil {
			return v * baseFontSize
		}
	}
	if strings.HasSuffix(value, "em") {
		num := strings.TrimSuffix(value, "em")
		v, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
		if err == nil {
			return v * baseFontSize
		}
	}
	if strings.HasSuffix(value, "px") {
		num := strings.TrimSuffix(value, "px")
		v, err := strconv.ParseFloat(strings.TrimSpace(num), 64)
		if err == nil {
			return v
		}
	}
	// Try plain number
	v, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return v
	}
	return 0
}

// cssFontWeight converts CSS font-weight to Figma fontStyle string.
func cssFontWeight(weight string) string {
	weight = strings.TrimSpace(weight)
	switch weight {
	case "100":
		return "Thin"
	case "200":
		return "ExtraLight"
	case "300":
		return "Light"
	case "400", "normal":
		return "Regular"
	case "500":
		return "Medium"
	case "600":
		return "SemiBold"
	case "700", "bold":
		return "Bold"
	case "800":
		return "ExtraBold"
	case "900":
		return "Black"
	default:
		return "Regular"
	}
}

// Helper functions

func parseFloatDef(s string, def float64) float64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return def
	}
	return v
}

func clampInt(v float64, min, max int) int {
	i := int(math.Round(v))
	if i < min {
		return min
	}
	if i > max {
		return max
	}
	return i
}
