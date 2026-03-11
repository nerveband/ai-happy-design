package tools

// GuideDepth returns deep-dive design intelligence on depth, shadows, and surface hierarchy.
// Teaches principles and valid ranges, not hardcoded recipes.
func GuideDepth() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Depth communicates hierarchy through subtle surface shifts. Choose ONE depth strategy and commit to it across the entire design. Mixing strategies creates visual incoherence.",
		"strategies": []map[string]interface{}{
			{"name": "borders-only", "style": "flat", "bestFor": "Technical tools, dense dashboards, data-heavy UIs. Clean and focused.", "technique": "1px strokes at 5-12% opacity. No shadows."},
			{"name": "subtle-shadows", "style": "soft lift", "bestFor": "Approachable apps, marketing sites, cards-based layouts.", "technique": "Single shadow per element. Low offset, moderate blur."},
			{"name": "layered-shadows", "style": "premium dimensional", "bestFor": "Luxury brands, landing pages, hero-driven designs.", "technique": "2-3 shadow layers per element. Increasing offset and blur at each layer."},
			{"name": "surface-tints", "style": "flat with depth", "bestFor": "Minimal designs, dark themes, glass-style UIs.", "technique": "Background color shifts (lighter = higher). No shadows needed."},
		},
		"shadows": map[string]interface{}{
			"principle": "Shadows communicate elevation. Higher elements cast larger, softer, more offset shadows. Lower elements cast tight, subtle shadows.",
			"scaling": map[string]interface{}{
				"formula":  "For level 1-5: offsetY = level * 2, blur = level * 4, opacity = 0.05 + level * 0.05",
				"level1":   "offsetY:2, blur:4, opacity:0.10 — subtle lift (cards, inputs)",
				"level2":   "offsetY:4, blur:8, opacity:0.15 — medium lift (dropdowns, popovers)",
				"level3":   "offsetY:6, blur:12, opacity:0.20 — elevated (modals, floating panels)",
				"level4":   "offsetY:8, blur:16, opacity:0.25 — high (dialogs, alerts)",
				"level5":   "offsetY:10, blur:20, opacity:0.30 — highest (full-screen overlays)",
			},
			"ranges": map[string]interface{}{
				"offsetY":  "1-16px. Always positive (light from above). offsetX usually 0.",
				"blur":     "4-48px. Larger blur = softer, more diffuse shadow.",
				"spread":   "0-4px typically. Negative spread for tight, inset-looking shadows.",
				"opacity":  "0.05-0.35. Lower = subtle, higher = dramatic. Most UI shadows: 0.08-0.20.",
			},
			"coloredShadows": "Use accent color at low opacity (0.10-0.20) for glow effects on CTAs or hero elements. Radius 16-32px.",
			"paramNames": "shadowType:'DROP_SHADOW', offsetX, offsetY, radius (blur), spread, color:{r,g,b,a}",
		},
		"glass": map[string]interface{}{
			"_when":     "Use when content overlaps complex backgrounds (images, gradients, other content).",
			"components": "fill (white/black at opacity 0.04-0.12) + BACKGROUND_BLUR (radius 8-24) + stroke (white/black at opacity 0.05-0.15).",
			"intensity": map[string]interface{}{
				"light":  "fill opacity 0.04, blur radius 8, stroke opacity 0.05. Barely there — subtle frosted.",
				"medium": "fill opacity 0.08, blur radius 16, stroke opacity 0.10. Visible glass effect.",
				"heavy":  "fill opacity 0.12, blur radius 24, stroke opacity 0.15. Strong frosted glass.",
			},
			"rule":      "Glass is an accent technique, not a default surface treatment. Overuse dulls the effect.",
			"paramName": "effect.apply_glass {nodeId, intensity:'light'|'medium'|'heavy'}",
		},
		"blur": map[string]interface{}{
			"layerBlur":      "LAYER_BLUR — blurs the element itself. Use for background depth (blur a shape behind content to soften it).",
			"backgroundBlur": "BACKGROUND_BLUR — blurs what's behind the element (frosted glass). The element must have a semi-transparent fill.",
			"rule":           "Never blur readable content. Blur is for decorative or background elements only.",
			"paramName":      "effect.add_blur {blurType:'LAYER_BLUR'|'BACKGROUND_BLUR', radius:number}",
		},
		"borders": map[string]interface{}{
			"cards":    "1px stroke at 5-12% opacity. Provides subtle separation without heaviness.",
			"inputs":   "1px stroke, slightly DARKER than surroundings. Inset feel communicates 'type here.'",
			"dividers": "1px line at 5-10% opacity. Horizontal rules between content sections.",
			"active":   "Increase stroke opacity or switch to accent color for focus/active borders.",
			"paramName": "strokeWidth (number), strokeAlign:'INSIDE'|'OUTSIDE'|'CENTER', color via paint.set_stroke",
		},
		"paramReference": map[string]interface{}{
			"offsetX":     "number — horizontal shadow offset",
			"offsetY":     "number — vertical shadow offset (positive = downward)",
			"radius":      "number — blur radius for shadows and blur effects",
			"spread":      "number — shadow spread (positive = larger, negative = tighter)",
			"color":       "object {r,g,b,a} — shadow/stroke color with opacity",
			"opacity":     "number 0.0-1.0 — element opacity",
			"blurType":    "string — 'LAYER_BLUR' or 'BACKGROUND_BLUR'",
			"strokeWidth": "number — border thickness in pixels",
		},
		"test": "Elevation test: does depth communicate hierarchy? Higher/more prominent elements should feel more important or urgent. If all surfaces look the same, depth isn't working.",
	}
}
