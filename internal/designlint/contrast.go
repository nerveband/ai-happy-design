package designlint

import (
	"math"
	"strconv"
	"strings"
)

// ContrastRatio computes WCAG 2.0 contrast ratio between two hex colors.
// Returns a value between 1 (no contrast) and 21 (max contrast).
func ContrastRatio(hex1, hex2 string) float64 {
	l1 := relativeLuminance(hex1)
	l2 := relativeLuminance(hex2)
	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)
	return (lighter + 0.05) / (darker + 0.05)
}

func relativeLuminance(hex string) float64 {
	r, g, b := parseHexRGB(hex)
	rr := linearize(float64(r) / 255.0)
	gg := linearize(float64(g) / 255.0)
	bb := linearize(float64(b) / 255.0)
	return 0.2126*rr + 0.7152*gg + 0.0722*bb
}

func linearize(v float64) float64 {
	if v <= 0.03928 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func parseHexRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	if len(hex) < 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 32)
	g, _ := strconv.ParseInt(hex[2:4], 16, 32)
	b, _ := strconv.ParseInt(hex[4:6], 16, 32)
	return int(r), int(g), int(b)
}

// AdjustForContrast returns a color adjusted to meet minimum contrast ratio against bg.
func AdjustForContrast(fgHex, bgHex string, minRatio float64) string {
	ratio := ContrastRatio(fgHex, bgHex)
	if ratio >= minRatio {
		return fgHex
	}

	bgLum := relativeLuminance(bgHex)
	if bgLum > 0.5 {
		return "#222222"
	}
	return "#FFFFFF"
}
