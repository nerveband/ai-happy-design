package tools

// GuideTypography returns deep-dive design intelligence on typography systems.
// Teaches principles and valid ranges, not hardcoded recipes.
func GuideTypography() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Typography is a system, not a font choice. Every text element exists in a hierarchy — size, weight, tracking, and opacity work together to communicate importance.",
		"hierarchy": map[string]interface{}{
			"levels": []map[string]interface{}{
				{"level": "primary", "role": "display/hero — the ONE focal point", "fontSize": "48-120px+", "fontStyle": "Bold/ExtraBold (700-800)", "letterSpacing": "tighten -0.02 to -0.03em", "color": "full opacity"},
				{"level": "secondary", "role": "heading/subheading — section anchors", "fontSize": "24-48px", "fontStyle": "SemiBold/Bold (600-700)", "letterSpacing": "-0.01 to -0.02em", "color": "full or 90% opacity"},
				{"level": "tertiary", "role": "body — readable content", "fontSize": "14-20px", "fontStyle": "Regular/Medium (400-500)", "letterSpacing": "0 (default)", "color": "65-80% opacity"},
				{"level": "quaternary", "role": "caption/label — supporting detail", "fontSize": "11-14px", "fontStyle": "Medium (500)", "letterSpacing": "0.01-0.02em", "color": "45-65% opacity"},
			},
			"rule": "Differentiate levels by size + weight + tracking + opacity — not size alone. Each level should be visually distinct at a glance.",
		},
		"letterSpacing": map[string]interface{}{
			"_rule":      "Letter-spacing compensates for optical effects at different sizes.",
			"allCaps":    "textCase:'UPPER' MANDATORY: 0.06-0.10em. ALL CAPS without tracking looks cramped.",
			"smallText":  "Below 13px: add 0.01-0.02em for readability.",
			"headings":   "32px+: tighten -0.01 to -0.02em. Large text has natural optical spacing.",
			"display":    "48px+: tighten -0.02 to -0.03em. Display type needs tight tracking for visual cohesion.",
			"body":       "14-20px: 0 (default). Body text reads best at natural spacing.",
			"paramName":  "letterSpacing (number, in em units relative to fontSize)",
		},
		"lineHeight": map[string]interface{}{
			"_critical": "ALWAYS pass lineHeightUnit:'PERCENT'. Without it, value is interpreted as PIXELS.",
			"body":      "1.5-1.7 (150-170%). Comfortable reading for paragraphs.",
			"headlines": "1.0-1.2 (100-120%). Tight for visual impact.",
			"display":   "0.9-1.1 (90-110%). Very tight — display text is meant to be seen, not read.",
			"uiText":    "1.0 (100%). Buttons, labels, badges — single line, no extra space.",
			"formula":   "nodeHeight = fontSize * (lineHeight / 100). A 152px font at 150% lineHeight = 228px node height.",
			"paramName": "lineHeight (number) + lineHeightUnit:'PERCENT'",
		},
		"lineLength": map[string]interface{}{
			"rule":      "Max 50-75 characters per line for body text. Beyond 75 chars, reading comprehension drops.",
			"formula":   "Approximate max-width: fontSize * 30 for proportional fonts.",
			"paramName": "width (on text.create — sets wrapping boundary)",
		},
		"colorSystem": map[string]interface{}{
			"primary":   "Near-black, never pure #000000. Use #0B0B0B or #1A1A1A. Full hierarchy visibility.",
			"secondary": "65% opacity of primary. Section labels, metadata.",
			"tertiary":  "45% opacity of primary. Timestamps, helper text.",
			"disabled":  "30% opacity of primary. Inactive states.",
			"paramName": "color (hex string '#RRGGBB' or {r,g,b,a} object on text.create)",
		},
		"fontWeight": map[string]interface{}{
			"principle": "Larger size = lighter weight is OK. Smaller size = heavier weight needed for legibility.",
			"maxWeights": "3-4 weights in a design. More creates visual noise.",
			"mapping":   "400=Regular, 500=Medium, 600=SemiBold, 700=Bold, 800=ExtraBold",
			"paramName": "fontStyle (string: 'Regular', 'Medium', 'SemiBold', 'Bold', 'ExtraBold')",
		},
		"constraints": map[string]interface{}{
			"maxFamilies": "2 font families. One for headings, one for body. More creates inconsistency.",
			"maxSizes":    "6-8 distinct font sizes. Define a scale and stick to it.",
			"maxWeights":  "3-4 weights total across the design.",
			"paramNames":  "fontSize (number), fontFamily (string), fontStyle (string)",
		},
		"paramReference": map[string]interface{}{
			"fontSize":      "number — size in pixels",
			"fontFamily":    "string — font family name (e.g., 'Inter', 'Space Grotesk')",
			"fontStyle":     "string — weight/style name (e.g., 'Regular', 'Bold', 'SemiBold')",
			"letterSpacing": "number — tracking in em units",
			"lineHeight":    "number — line height value (use with lineHeightUnit:'PERCENT')",
			"textCase":      "string — 'UPPER', 'LOWER', 'TITLE', 'ORIGINAL'",
			"color":         "string or object — hex '#RRGGBB' or {r,g,b,a}",
		},
		"test": "Squint test: blur or squint at the design — are 3+ hierarchy levels still visually distinct? If everything looks the same size, the type system is too flat.",
	}
}
