package tools

// GuideLayout returns deep-dive design intelligence on layout and spacing systems.
// Teaches principles and valid ranges, not hardcoded recipes.
func GuideLayout() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Spacing communicates relationships. Elements close together are related; elements far apart are separate. Every pixel of spacing is a design decision.",
		"grid": map[string]interface{}{
			"rule":    "Every x, y, width, height, padding, and itemSpacing value should be a multiple of 8px. Use 4px ONLY for tight internal spacing (icon-to-label gaps).",
			"values":  "Common: 8, 16, 24, 32, 48, 64, 96, 128. These are your vocabulary — use them consistently.",
			"formula": "Round to grid: Math.round(value / 8) * 8. For 4px sub-grid: Math.round(value / 4) * 4.",
		},
		"proximity": map[string]interface{}{
			"_rule":          "Proximity = relationship. Closer means 'these belong together.'",
			"withinGroup":    "8-16px. Elements that form a unit (icon + label, title + subtitle).",
			"betweenGroups":  "24-48px. Distinct groups within a section (card content blocks, form field groups).",
			"betweenSections": "48-96px. Major content sections (hero → features → footer).",
			"formula":        "sectionGap >= 2 * groupGap >= 4 * elementGap. Maintain clear ratio between levels.",
		},
		"composition": map[string]interface{}{
			"pattern":  "Root frame → background layer → content wrapper (auto-layout) → children. This 3-layer structure handles 90% of layouts.",
			"rootFrame": "FIXED width and height. This is your canvas — explicit dimensions.",
			"background": "Full-bleed fill, gradient, or image. Sent to back with layer.send_to_back.",
			"content":   "Auto-layout frame (VERTICAL or HORIZONTAL) centered in root. This is where all content lives.",
		},
		"autoLayout": map[string]interface{}{
			"_rule":      "Auto-layout IS CSS flexbox. If you know flexbox, you know auto-layout.",
			"direction":  "layoutMode: 'VERTICAL' (flex-direction:column) or 'HORIZONTAL' (flex-direction:row).",
			"gap":        "itemSpacing: number — equivalent to CSS gap.",
			"padding":    "padding: number (uniform) or paddingTop/Right/Bottom/Left (per-side).",
			"mainAxis":   "primaryAxisAlignItems: 'MIN' (flex-start), 'CENTER', 'MAX' (flex-end), 'SPACE_BETWEEN'.",
			"crossAxis":  "counterAxisAlignItems: 'MIN', 'CENTER', 'MAX'.",
			"centering":  "primaryAxisAlign:'CENTER' + counterAxisAlign:'CENTER' = both-axis centering.",
			"paramNames": "layoutMode, itemSpacing, padding, primaryAxisAlign, counterAxisAlign",
		},
		"sizing": map[string]interface{}{
			"HUG":   "Shrink-wraps to content. Use for: buttons, badges, nav items, cards with variable content. Prevents empty whitespace.",
			"FILL":  "Stretches to fill parent. ONLY valid inside auto-layout parent. Use for: full-width text, content areas, section backgrounds.",
			"FIXED": "Exact pixel dimensions. Use for: root frames, images, icons. NOT for text containers (they should HUG or FILL).",
			"rule":  "Default to HUG. Use FILL for stretching within auto-layout. FIXED only for root frames and known-dimension elements.",
		},
		"balance": map[string]interface{}{
			"_rule":    "Sibling elements MUST match. Unbalanced siblings look broken, even if individually well-designed.",
			"match":    "Same height, padding, itemSpacing, cornerRadius, and text sizes for sibling cards/sections.",
			"formula":  "cardWidth = (rowWidth - (N-1) * gap) / N. All sibling cards get the SAME computed width.",
			"textMatch": "Text at the same hierarchy level across siblings: same fontSize, fontStyle, color.",
		},
		"asymmetry": map[string]interface{}{
			"rule": "Intentional asymmetry in ONE element creates visual tension and focal points. Everything asymmetric = chaos. Break the grid deliberately and for a reason.",
		},
		"proportions": map[string]interface{}{
			"rule":    "Proportions declare importance. Larger = more important. The hero section should visually dominate. Cards should be notably smaller than the hero but larger than captions.",
			"minSizes": "Card height >= 15-20% of canvas height. Content area >= 60% of canvas. Hero text >= 2x body text.",
		},
		"paramReference": map[string]interface{}{
			"x":                 "number — horizontal position relative to parent",
			"y":                 "number — vertical position relative to parent",
			"width":             "number — element width in pixels",
			"height":            "number — element height in pixels",
			"padding":           "number — uniform padding (or paddingTop/Right/Bottom/Left)",
			"itemSpacing":       "number — gap between auto-layout children",
			"layoutMode":        "string — 'VERTICAL' or 'HORIZONTAL'",
			"primaryAxisAlign":  "string — 'MIN', 'CENTER', 'MAX', 'SPACE_BETWEEN'",
			"counterAxisAlign":  "string — 'MIN', 'CENTER', 'MAX'",
		},
		"test": "Proximity test: can you tell which elements are grouped together just by spacing? If groups aren't obvious at a glance, spacing ratios need more contrast.",
	}
}
