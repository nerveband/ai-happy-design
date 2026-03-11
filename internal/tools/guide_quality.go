package tools

// GuideQuality returns a quality gate checklist for design self-verification.
// Each check is actionable and specific — pass or fail, no ambiguity.
func GuideQuality() map[string]interface{} {
	return map[string]interface{}{
		"_philosophy": "Quality is not a final step — it is the standard at every step. These gates catch the difference between 'technically correct' and 'professionally designed.' Run them before declaring any design complete.",
		"gates": []map[string]interface{}{
			{
				"name":     "intention",
				"question": "Can you justify every color, font, and spacing choice for THIS specific project?",
				"fail":     "Choices feel generic or borrowed from templates. No clear connection between visual decisions and project purpose.",
				"action":   "Revisit color (guide --topic color) and typography (guide --topic typography). Make each choice deliberate.",
			},
			{
				"name":     "swap",
				"question": "If you changed the accent color and font, would this design feel meaningfully different?",
				"fail":     "Design looks the same with any color/font combo. Visual identity is absent.",
				"action":   "The accent and typography should reinforce the project's personality. Redesign with purpose.",
			},
			{
				"name":     "squint",
				"question": "Are 3+ hierarchy levels visible when you blur or squint at the design?",
				"fail":     "Everything looks the same size/weight. No focal point. No visual flow.",
				"action":   "Increase contrast between hierarchy levels. Bigger size jumps, bolder weights, more opacity difference.",
			},
			{
				"name":     "contrast",
				"question": "Does all text meet minimum contrast ratios?",
				"criteria": "4.5:1 for body text (<=16px). 3:1 for large text (18px+ bold, 24px+). 3:1 for UI components (borders, icons, controls).",
				"action":   "Darken text or lighten background. Use contrast checker or compute relative luminance.",
			},
			{
				"name":     "gridAlignment",
				"question": "Are all x, y, width, height, padding, and spacing values on the 8px grid?",
				"fail":     "Values like 13, 27, 45, 97 — not multiples of 8 (or 4 for tight spacing).",
				"action":   "Round all values: Math.round(value / 8) * 8. Use 4px sub-grid only for icon-label gaps.",
			},
			{
				"name":     "balance",
				"question": "Do sibling elements match in dimensions, padding, spacing, radius, and text sizes?",
				"fail":     "Cards at different heights. Inconsistent padding. Mixed corner radii among peers.",
				"action":   "Audit sibling groups. All cards: same height, padding, itemSpacing, cornerRadius, text sizing.",
			},
			{
				"name":     "stateCoverage",
				"question": "Do interactive elements have hover/active/focus/disabled states? Do data views have loading/empty/error?",
				"fail":     "Buttons with no hover. Lists with no empty state. Forms with no error state.",
				"action":   "Design all 5 interactive states and all 3 data states. See guide --topic states.",
			},
			{
				"name":     "naming",
				"question": "Is every layer named by its role? No defaults like 'Frame 47' or 'Rectangle 12'?",
				"fail":     "Default Figma names anywhere in the layer tree.",
				"action":   "Name every layer: 'Hero Background', 'Card - Feature', 'CTA Button', 'Section Divider'.",
			},
			{
				"name":     "typographySystem",
				"question": "Max 2 font families, 6-8 font sizes, 3-4 weights? Proper letter-spacing on ALL CAPS text?",
				"fail":     "More than 2 families. Random font sizes. ALL CAPS without tracking.",
				"action":   "Consolidate to a type scale. Add letterSpacing 0.06-0.10em to all textCase:'UPPER' elements.",
			},
			{
				"name":     "domainFit",
				"question": "Does this look like it belongs to THIS product's world, or could it be any product?",
				"fail":     "Generic startup look. No visual personality. Interchangeable with any SaaS template.",
				"action":   "Identify 2-3 visual traits unique to this project and weave them consistently throughout.",
			},
			{
				"name":     "exportVerification",
				"question": "Have you exported at scale 2 and visually verified the result?",
				"fail":     "Reported success without visual verification. Overlaps, overflow, or alignment issues only visible in export.",
				"action":   "export.image {nodeId, scale:2}. Open the exported file and review. Fix any issues found.",
			},
		},
		"workflow": "Run these gates in order. Fix failures before proceeding. Gates 1-3 (intention/swap/squint) catch design direction problems early. Gates 4-9 catch execution problems. Gates 10-11 are final verification.",
		"test": "Final question: would you be proud to show this to a design-savvy colleague? If hesitant, identify which gate is failing and fix it.",
	}
}
