package tools

// GuideColor returns deep-dive design intelligence on color systems.
// Teaches principles and valid ranges, not hardcoded recipes.
func GuideColor() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Color is restraint. Every color choice must be intentional and justifiable — 'Why this color for THIS project?' not 'it looks nice.' One purposeful accent beats five competing colors.",
		"palette": map[string]interface{}{
			"layers": []map[string]interface{}{
				{"layer": "neutrals", "coverage": "70-90% of UI", "role": "Backgrounds, surfaces, text, borders. The foundation everything else rests on."},
				{"layer": "primary accent", "coverage": "5-10%", "role": "Brand identity, CTAs, active states, key highlights. ONE color used consistently."},
				{"layer": "semantic", "coverage": "status only", "role": "Success (green), warning (amber), danger (red). Each needs base + tint + border + on-color."},
				{"layer": "effects", "coverage": "rare", "role": "Gradients, glows, decorative. Used sparingly for emphasis."},
			},
			"rule": "Start with neutrals. Add ONE accent. Add semantic only where needed. Effects last and sparingly.",
		},
		"neutrals": map[string]interface{}{
			"scale":    "10-12 steps from near-white (#FAFAFA) to near-black (#0A0A0A).",
			"noPure":   "Never pure #000000 on pure #FFFFFF. Use #0B0B0B on #FAFAFA or similar. Pure black/white creates harsh vibration.",
			"naming":   "Name by purpose: 'surface', 'surface-elevated', 'border', 'text-primary' — not 'gray-100', 'gray-200'.",
			"paramName": "color (hex string or {r,g,b,a})",
		},
		"accent": map[string]interface{}{
			"rule":          "One accent color, used purposefully. Max saturation ~80% — fully saturated colors feel aggressive.",
			"consistency":   "Same accent for: primary buttons, active nav, links, focus rings, key highlights.",
			"tints":         "Derive tints by reducing opacity or mixing with background. accent-subtle = 10% opacity on surface.",
			"intentionTest": "If you swapped the accent for blue, would this design feel meaningfully different? If not, the accent is decorative, not intentional.",
		},
		"semantic": map[string]interface{}{
			"success": "Green family — base, tint (10% on surface), border (20% on surface), on-color (white or dark text).",
			"warning": "Amber family — same 4-token structure. Amber, not yellow (yellow has poor contrast).",
			"danger":  "Red family — same 4-token structure. Use for destructive actions and error states.",
			"rule":    "Semantic colors are functional, not decorative. Never use red for a heading just because it stands out.",
		},
		"darkTheme": map[string]interface{}{
			"_rule":       "Dark theme is a SEPARATE system, not colors inverted. Design it independently.",
			"background":  "#0F0F0F (not #000000). Pure black is a void — subtle warmth feels intentional.",
			"text":        "#F0F0F0 (not #FFFFFF). Slightly off-white reduces eye strain.",
			"elevation":   "Higher surfaces = lighter. Background → surface (+2% white) → elevated (+4% white) → overlay (+6% white).",
			"borders":     "White at 5-12% opacity. Subtle enough to separate, not enough to compete.",
			"accent":      "Same hue as light theme but may need saturation/lightness adjustment for dark backgrounds.",
		},
		"gradients": map[string]interface{}{
			"purpose":    "Gradients serve brand or hierarchy. Hero backgrounds, CTAs, accent overlays — not every surface.",
			"readability": "Never place text directly on a gradient without sufficient contrast at ALL points.",
			"subtlety":   "Subtle > dramatic. 2-3 stops maximum. Stops should be close in hue for harmony.",
			"opacity":    "For overlay gradients, keep stop opacity low (0.05-0.15). Dramatic gradients are for hero backgrounds only.",
			"paramName":  "stops: [{position:0, color:'#hex'}, {position:1, color:'#hex'}], gradientType: 'LINEAR'|'RADIAL'",
		},
		"contrast": map[string]interface{}{
			"bodyText":    "4.5:1 minimum (WCAG AA). Body text <=16px against its background.",
			"largeText":   "3:1 minimum. Large text = 18px+ bold OR 24px+ regular.",
			"uiComponents": "3:1 minimum. Borders, icons, form controls against their background.",
			"formula":     "Use relative luminance: L = 0.2126*R + 0.7152*G + 0.0722*B. Ratio = (L1 + 0.05) / (L2 + 0.05).",
		},
		"paramReference": map[string]interface{}{
			"color": "string '#RRGGBB' or object {r,g,b,a} where r/g/b are 0.0-1.0 floats, a is optional opacity 0.0-1.0",
		},
		"test": "Swap test: if you changed the accent color to blue, would this design feel meaningfully different? If not, the color choices are arbitrary — go back and make them intentional for THIS project.",
	}
}
