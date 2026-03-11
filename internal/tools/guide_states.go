package tools

// GuideStates returns deep-dive design intelligence on interactive states and data view states.
// Teaches principles and valid ranges, not hardcoded recipes.
func GuideStates() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Missing states feel broken — even if technically functional. Every interactive element and every data view needs explicit state design. Users judge quality by how things respond, not just how they look.",
		"interactiveStates": map[string]interface{}{
			"_rule": "Every interactive element (button, link, input, toggle, card) needs all 5 states designed.",
			"states": []map[string]interface{}{
				{"state": "default", "description": "Resting appearance. Must clearly communicate 'this is interactive' through color, shape, or visual weight."},
				{"state": "hover", "description": "Subtle shift on mouse-over. Color/opacity change, 150ms transition. NOT scale transforms that shift layout."},
				{"state": "active", "description": "Pressed/clicking. -1px translateY or scale(0.98) for tactile 'press' feel. Brief, snappy response."},
				{"state": "focus", "description": "Keyboard navigation. Visible ring (2px, accent color or high-contrast) for accessibility. NEVER remove focus indicators."},
				{"state": "disabled", "description": "Inactive/unavailable. 50% opacity, no pointer events. Must still be recognizable as the same element type."},
			},
			"hover": map[string]interface{}{
				"technique":  "Lighten or darken background by 5-10%. Or shift opacity by 0.1. Subtle, not dramatic.",
				"timing":     "150ms transition. Fast enough to feel responsive, slow enough to be smooth.",
				"avoid":      "Scale transforms that move nearby elements. Color changes so dramatic they look like a different element.",
			},
			"active": map[string]interface{}{
				"technique": "translateY(-1px) or scale(0.98) for physical press feel. Or darken background 10-15%.",
				"timing":    "Immediate (no transition delay). The press should feel instant.",
			},
			"focus": map[string]interface{}{
				"technique": "2px outline or ring, 2px offset, accent color or high-contrast color.",
				"rule":      "Focus indicators are non-negotiable for accessibility. Design them to look intentional, not like a browser default.",
			},
		},
		"dataViewStates": map[string]interface{}{
			"_rule": "Every data view (list, table, grid, feed, detail) needs all 3 states designed.",
			"states": []map[string]interface{}{
				{"state": "loading", "description": "Skeleton screens matching content shape. NOT spinners. Skeletons feel faster and prevent layout shift."},
				{"state": "empty", "description": "Explain WHY empty + WHAT to do next. Icon + headline + description + primary action button."},
				{"state": "error", "description": "Specific error message (not 'Something went wrong') + clear next action (retry button, support link)."},
			},
			"loading": map[string]interface{}{
				"technique": "Skeleton shapes matching the expected content layout. Animated pulse (opacity 0.05 to 0.12, 1.5s loop).",
				"rule":      "Skeleton bones should match the approximate size and position of real content. This prevents layout shift when data loads.",
			},
			"empty": map[string]interface{}{
				"structure": "Centered in container: illustration/icon (optional) → headline ('No items yet') → description ('Add your first item to get started') → action button.",
				"rule":      "Empty states are an opportunity, not a dead end. Guide the user to the next action.",
			},
			"error": map[string]interface{}{
				"structure": "Icon (warning/error) → specific message → action button ('Retry', 'Go back', 'Contact support').",
				"rule":      "Be specific about what went wrong. 'Failed to load messages' not 'Something went wrong.' Always offer a next step.",
			},
		},
		"animationTiming": map[string]interface{}{
			"micro":       "150ms — hover, toggle, checkbox, small state changes.",
			"reveal":      "200-250ms — dropdowns, tooltips, toasts appearing.",
			"major":       "300-500ms — page transitions, modal open/close, large content reveals.",
			"easing": map[string]interface{}{
				"entrance": "Deceleration (ease-out). Fast start, gentle landing. Objects arriving into view.",
				"exit":     "Acceleration (ease-in). Gentle start, fast departure. Objects leaving view.",
				"move":     "Ease-in-out. Smooth both ways. Objects moving within view.",
				"avoid":    "Linear easing for UI motion. Linear feels mechanical and unnatural.",
			},
		},
		"touchTargets": map[string]interface{}{
			"minimum":  "44x44px minimum touch target. Even if the visual element is smaller, the hit area must be at least 44x44.",
			"spacing":  "8px minimum between adjacent touch targets to prevent mis-taps.",
			"paramName": "width and height on interactive elements (buttons, toggles, links in mobile contexts)",
		},
		"paramReference": map[string]interface{}{
			"opacity":      "number 0.0-1.0 — for hover/disabled state opacity shifts",
			"color":        "string or {r,g,b,a} — for state color changes",
			"visible":      "boolean — for showing/hiding state variants",
			"width":        "number — minimum 44px for touch targets",
			"height":       "number — minimum 44px for touch targets",
		},
		"test": "Interaction audit: click/hover every interactive element — does each respond predictably? Check every data view: what does it look like loading, empty, and errored?",
	}
}
