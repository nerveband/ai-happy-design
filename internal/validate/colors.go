package validate

import "strings"

// cssColors maps CSS named colors to hex values.
var cssColors = map[string]string{
	"black":       "#000000", "white":       "#FFFFFF", "red":         "#FF0000",
	"green":       "#008000", "blue":        "#0000FF", "yellow":      "#FFFF00",
	"cyan":        "#00FFFF", "magenta":     "#FF00FF", "silver":      "#C0C0C0",
	"gray":        "#808080", "grey":        "#808080", "maroon":      "#800000",
	"olive":       "#808000", "lime":        "#00FF00", "aqua":        "#00FFFF",
	"teal":        "#008080", "navy":        "#000080", "fuchsia":     "#FF00FF",
	"purple":      "#800080", "orange":      "#FFA500", "pink":        "#FFC0CB",
	"brown":       "#A52A2A", "coral":       "#FF7F50", "crimson":     "#DC143C",
	"gold":        "#FFD700", "indigo":      "#4B0082", "ivory":       "#FFFFF0",
	"khaki":       "#F0E68C", "lavender":    "#E6E6FA", "linen":       "#FAF0E6",
	"orchid":      "#DA70D6", "peru":        "#CD853F", "plum":        "#DDA0DD",
	"salmon":      "#FA8072", "sienna":      "#A0522D", "tan":         "#D2B48C",
	"thistle":     "#D8BFD8", "tomato":      "#FF6347", "turquoise":   "#40E0D0",
	"violet":      "#EE82EE", "wheat":       "#F5DEB3", "transparent": "#00000000",
}

// ResolveNamedColor converts a CSS named color to hex. Returns "" if not found.
func ResolveNamedColor(name string) string {
	return cssColors[strings.ToLower(strings.TrimSpace(name))]
}
