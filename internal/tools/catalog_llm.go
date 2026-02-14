package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// LLMCatalog returns an enriched machine-readable catalog with examples and
// recommended execution patterns for LLM agents.
func LLMCatalog() map[string]interface{} {
	toolNames := make([]string, 0, len(toolDescriptions))
	for name := range toolDescriptions {
		toolNames = append(toolNames, name)
	}
	sort.Strings(toolNames)

	toolsOut := make(map[string]interface{}, len(toolNames))
	for _, toolName := range toolNames {
		actions := toolDescriptions[toolName]
		actionNames := make([]string, 0, len(actions))
		for action := range actions {
			actionNames = append(actionNames, action)
		}
		sort.Strings(actionNames)

		actionOut := make(map[string]interface{}, len(actionNames))
		for _, actionName := range actionNames {
			desc := actions[actionName]
			actionOut[actionName] = buildActionSpec(toolName, actionName, desc)
		}

		toolsOut[toolName] = map[string]interface{}{
			"actions": actionOut,
		}
	}

	return map[string]interface{}{
		"version": "3.1",
		"discovery": map[string]interface{}{
			"cliCatalog":       "ai-happy-design tools --json",
			"llmCatalog":       "ai-happy-design tools --llm --json",
			"mcpDescribe":      "describe(action='catalog')",
			"cliActions":       "ai-happy-design actions [domain]",
			"domainActionHint": "Prefer domain.action names like paint.set_solid over legacy command aliases.",
		},
		"workflow": map[string]interface{}{
			"rule": "CREATING = batch (one payload, many steps). EDITING = single command (precise, targeted).",
			"create": map[string]interface{}{
				"when":    "Building new designs, layouts, multi-element compositions, or anything with 3+ elements.",
				"how":     "Build a JSON array of operations with named steps and ${{steps.name.result.id}} interpolation. Send as one batch: 'ai-happy-design batch -f payload.json' or bulk.execute via MCP.",
				"why":     "One WebSocket connection, no process overhead per step, automatic ID chaining. 150 steps in ~6 seconds vs ~8 minutes with individual commands.",
				"example": "See designPatterns below for batch payload templates.",
			},
			"edit": map[string]interface{}{
				"when": "Changing a color, moving a node, resizing, renaming, or any single-property change on an existing node.",
				"how":  "Use a single command: 'ai-happy-design command paint.set_solid -p {\"nodeId\":\"1:2\",\"color\":\"#FF0000\"}'",
				"why":  "Fast, precise, no batch overhead needed for one operation.",
			},
		},
		"designThinking": map[string]interface{}{
			"_overview": "Think like a frontend developer. You already know HTML/CSS — use that knowledge. Mentally draft the design as HTML/CSS first, then translate each CSS property to Figma commands. This produces dramatically better designs than generating Figma commands directly.",
			"process": []string{
				"1. VISUALIZE: Picture the design as a webpage. What HTML elements would you use? What CSS would you write?",
				"2. STRUCTURE: Map HTML elements to Figma frames. <div> = frame, <p> = text, <img> = image fill, <section> = auto-layout frame.",
				"3. STYLE: Translate CSS properties to Figma commands using the cssToFigma map below.",
				"4. LAYER: Name every element descriptively. Group related elements in frames. Order layers front-to-back.",
				"5. VERIFY: Export and review. Does it match what the CSS would render?",
			},
			"cssToFigma": map[string]interface{}{
				"_hint":                           "When you think 'I would use this CSS', translate it to the corresponding Figma command.",
				"display: flex":                   "layout.set_auto_layout {direction, itemSpacing, padding}",
				"flex-direction: column":           "layout.set_auto_layout {direction: VERTICAL}",
				"flex-direction: row":              "layout.set_auto_layout {direction: HORIZONTAL}",
				"justify-content: center":          "layout.set_alignment {primaryAxisAlign: CENTER}",
				"justify-content: space-between":   "layout.set_alignment {primaryAxisAlign: SPACE_BETWEEN}",
				"align-items: center":              "layout.set_alignment {counterAxisAlign: CENTER}",
				"gap: 16px":                        "layout.set_spacing {itemSpacing: 16}",
				"padding: 24px":                    "layout.set_padding {padding: 24}",
				"padding: 16px 24px":               "layout.set_padding {paddingTop:16, paddingBottom:16, paddingLeft:24, paddingRight:24}",
				"background-color":                 "paint.set_solid {color, opacity}",
				"background: linear-gradient(...)": "paint.set_gradient {type:LINEAR, stops:[{position:0,color},{position:1,color}]}",
				"background: radial-gradient(...)": "paint.set_gradient {type:RADIAL, stops:[...]}",
				"border: 1px solid color":          "paint.set_stroke {color, width:1, strokeAlign:INSIDE}",
				"border-radius: 16px":              "node.set_corner_radius {radius:16}",
				"box-shadow: 0 4px 16px rgba(...)": "effect.add_shadow {shadowType:DROP_SHADOW, offsetX:0, offsetY:4, radius:16, color}",
				"box-shadow: inset ...":            "effect.add_shadow {shadowType:INNER_SHADOW, ...}",
				"filter: blur(8px)":                "effect.add_blur {blurType:LAYER_BLUR, radius:8}",
				"backdrop-filter: blur(16px)":      "effect.add_blur {blurType:BACKGROUND_BLUR, radius:16}",
				"opacity: 0.5":                     "node.set_opacity {opacity:0.5}",
				"mix-blend-mode: multiply":         "node.set_blend_mode {blendMode:MULTIPLY}",
				"color: white":                     "text.set_color {color:{r:1,g:1,b:1}}",
				"font-weight: 700":                 "text.set_weight {weight:700}",
				"font-size: 24px":                  "text.set_size {fontSize:24}",
				"text-transform: uppercase":        "text.set_case {textCase:UPPER}",
				"letter-spacing: 2px":              "text.set_letter_spacing {letterSpacing:2}",
				"line-height: 1.5":                 "text.set_line_height {value:150, unit:PERCENT}",
				"text-decoration: underline":       "text.set_decoration {decoration:UNDERLINE}",
				"z-index (higher)":                 "layer.bring_to_front",
				"z-index (lower)":                  "layer.send_to_back",
				"position: absolute":               "Create element in parent frame, set x/y manually (outside auto-layout flow)",
				"overflow: hidden":                 "Frames clip content by default. Use node.set_corner_radius on the frame for rounded clipping.",
				"width: 100%":                      "layout.set_sizing {counterAxisSizing:FIXED} + match parent width, or use FILL mode in auto-layout",
			},
			"visualHierarchy": map[string]interface{}{
				"_rule": "Every design has a visual hierarchy. Elements compete for attention. Use size, weight, color, contrast, and effects to control what the viewer sees first, second, and third.",
				"levels": map[string]interface{}{
					"primary": map[string]interface{}{
						"what":    "The ONE thing the viewer should see first. Hero title, main CTA, key number.",
						"how":     "Largest font size (hero scale), boldest weight (800), accent color, optional shadow glow.",
						"example": "fontSize:72, weight:800, color:accent. Or a big number like '106' in accent color.",
					},
					"secondary": map[string]interface{}{
						"what":    "Supporting elements. Subtitles, section headers, card titles, secondary CTAs.",
						"how":     "Medium font size (heading scale), semi-bold (600-700), white or slightly muted color.",
						"example": "fontSize:24-36, weight:600-700, color:white or warm accent.",
					},
					"tertiary": map[string]interface{}{
						"what":    "Details. Descriptions, body text, metadata, helper text.",
						"how":     "Smaller font size (body scale), regular weight (400-500), muted color (gray/low opacity).",
						"example": "fontSize:18-24, weight:400-500, color:{r:0.6,g:0.6,b:0.6}.",
					},
					"ambient": map[string]interface{}{
						"what":    "Background texture. Decorative shapes, subtle lines, particles, watermarks.",
						"how":     "Very low opacity (0.05-0.15), no text weight, accent or white color. Should NOT compete with content.",
						"example": "Ellipse with strokeColor at alpha 0.1, or rectangle fill at alpha 0.04.",
					},
				},
				"contrastRule": "Primary elements need HIGH contrast against background. Secondary = medium. Tertiary = low. Ambient = almost invisible. If everything is bold, nothing is bold.",
			},
			"designDecisions": map[string]interface{}{
				"_rule": "Match the visual treatment to the element's role. Think: what CSS would I write for this HTML element?",
				"whenToUseShadow": map[string]interface{}{
					"use": []string{
						"Cards and elevated surfaces — they 'float' above the background.",
						"Buttons and CTAs — subtle shadow adds depth and clickability.",
						"Modals and overlays — strong shadow separates from page.",
						"Images and avatars — subtle shadow adds polish.",
					},
					"avoid": []string{
						"Flat text — shadows on text are rarely good.",
						"Background decorations — they're behind everything.",
						"Borders and dividers — these are already subtle.",
					},
					"recipe_subtle":  "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.1}, offsetX:0, offsetY:2, radius:8}",
					"recipe_card":    "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.2}, offsetX:0, offsetY:4, radius:16}",
					"recipe_elevated": "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.3}, offsetX:0, offsetY:8, radius:32}",
					"recipe_glow":    "effect.add_shadow {shadowType:DROP_SHADOW, color:{accent_r,accent_g,accent_b,a:0.2}, offsetX:0, offsetY:4, radius:24}",
				},
				"whenToUseGradient": map[string]interface{}{
					"use": []string{
						"Hero backgrounds — gradient adds depth and visual interest vs flat color.",
						"CTAs and primary buttons — gradient makes them pop and feel premium.",
						"Accent overlays — subtle gradient overlay on images adds mood.",
						"Section dividers — gradient fade is more elegant than a hard line.",
						"Icon/badge backgrounds — small gradient adds richness to focal elements.",
					},
					"avoid": []string{
						"Body text — never put gradient on readable text (except hero titles as a stylistic choice).",
						"Every surface — overusing gradients makes designs look dated.",
						"Subtle cards — solid color with low opacity is usually better for card backgrounds.",
					},
					"recipe_bgSweep":  "paint.set_gradient {type:LINEAR, stops:[{position:0,color:dark},{position:1,color:slightlyLighter}]}",
					"recipe_accent":   "paint.set_gradient {type:LINEAR, stops:[{position:0,color:accent},{position:1,color:accentVariant}]}",
					"recipe_fade":     "paint.set_gradient {type:LINEAR, stops:[{position:0,color:{r:1,g:1,b:1,a:0.1}},{position:1,color:{r:1,g:1,b:1,a:0}}]}",
				},
				"whenToUseBlur": map[string]interface{}{
					"use": []string{
						"Glass/frosted effects — BACKGROUND_BLUR behind semi-transparent surface.",
						"Depth-of-field on background images — LAYER_BLUR on background layer.",
						"Focus attention — blur surrounding content to focus on modal/overlay.",
					},
					"avoid": []string{
						"Primary content — never blur things users need to read.",
						"Everything — blur is expensive. Use sparingly for premium feel.",
					},
					"recipe_glass": "paint.set_solid {color:{r:1,g:1,b:1,a:0.08}} + effect.add_blur {blurType:BACKGROUND_BLUR, radius:16} + paint.set_stroke {color:{r:1,g:1,b:1,a:0.1}, width:1}",
				},
				"whenToUseStroke": map[string]interface{}{
					"use": []string{
						"Card borders — subtle stroke defines edges on dark backgrounds.",
						"Input fields — strokes indicate interactive boundaries.",
						"Decorative outlines — stroke-only shapes (no fill) for ambient decoration.",
						"Dividers — 1px line between sections.",
						"Badges and tags — stroke adds definition to small elements.",
					},
					"recipe_subtle": "paint.set_stroke {color:{r:1,g:1,b:1,a:0.07}, width:1}",
					"recipe_accent": "paint.set_stroke {color:{accent,a:0.12}, width:1}",
					"recipe_decoRing": "paint.set_stroke {color:{accent,a:0.1}, width:2} + paint.remove_fill (stroke-only)",
				},
			},
			"layerOrganization": map[string]interface{}{
				"naming": map[string]interface{}{
					"rule": "Name every element by its content or role. Never leave default names like 'Frame 47' or 'Rectangle 12'.",
					"pattern": "[Type] - [Purpose]",
					"examples": map[string]interface{}{
						"frames":     "Hero Section, Card Grid, Card - Feature Name, Footer, Navigation Bar",
						"text":       "Hero Title, Card Title, Card Description, Footer Link, Badge Label",
						"shapes":     "Deco Ring, Background Gradient, Divider, Dot Indicator, Accent Line",
						"images":     "Avatar - User, Hero Background, Logo, Product Thumbnail",
					},
				},
				"grouping": map[string]interface{}{
					"rule": "Group related elements in frames. Every logical section should be a frame, not loose elements.",
					"structure": []string{
						"Root Frame (the canvas, e.g., 1080x1080)",
						"  ├── Background elements (decorative shapes, gradients) — BACK",
						"  ├── Content Frame (auto-layout, holds all structured content) — MIDDLE",
						"  │   ├── Header / Badge area",
						"  │   ├── Hero section (title, subtitle)",
						"  │   ├── Body section (cards, content grid)",
						"  │   └── Footer section (links, metadata)",
						"  └── Overlay elements (floating badges, stripes) — FRONT",
					},
				},
				"layerOrder": map[string]interface{}{
					"rule":    "Layer order = z-index. Back layers render first (behind), front layers render last (on top).",
					"pattern": "1. Background fills/gradients (BACK). 2. Decorative shapes (rings, particles). 3. Content frame (auto-layout). 4. Foreground overlays (stripes, badges). Use layer.send_to_back and layer.bring_to_front.",
				},
			},
			"componentRelationships": map[string]interface{}{
				"_rule": "Elements don't exist in isolation. Understand what CONTAINS what, what RELATES to what, and what EMPHASIZES what.",
				"containment": map[string]interface{}{
					"rule":     "A parent frame OWNS its children. Children inherit clipping from parent. Use parentId when creating children, or layer.move_to_parent after.",
					"examples": "A card frame contains its title, description, and icon. A nav bar contains logo, links, and button.",
				},
				"association": map[string]interface{}{
					"rule":     "Related elements share visual properties: same spacing, aligned edges, consistent colors within a group.",
					"examples": "All cards in a grid share the same radius, padding, font sizes, and spacing. Footer links share the same text color and size.",
				},
				"emphasis": map[string]interface{}{
					"rule":     "The most important element in a group should be visually distinct from its siblings. Use size, color, weight, or effects.",
					"examples": "In a card, the big number (106) is accent-colored and 48px while the label is white 24px. In a header, the brand name is accent-colored while 'Happy Design' is white.",
				},
				"contrast": map[string]interface{}{
					"rule":     "Adjacent elements need enough contrast to be distinguishable. Card vs background, text vs card, primary vs secondary.",
					"examples": "Dark card ({r:1,g:1,b:1,a:0.04}) on dark background ({r:0.1,g:0.1,b:0.1}) — the 4% white creates just enough lift. Add a 7% white stroke for edge definition.",
				},
				"proximity": map[string]interface{}{
					"rule":     "Elements that are close together are perceived as related. Use spacing to create groups. 8-16px within a group, 24-48px between groups.",
					"examples": "Card title and description: 8px apart (tight, same group). Card grid and footer: 32px apart (separate sections).",
				},
			},
		},
		"playbook": []string{
			"1. DISCOVER: Get catalog (describe action=catalog) before building commands.",
			"2. PROBE: document.get_info and document.get_selection to understand current state. Check existing frame positions.",
			"3. ASPECT RATIO: heightRatio = height / width. Pick layout strategy from typography.layoutByAspectRatio. Stories (1.78) = bigger text, taller cards, more spacing. Square (1.0) = compact. Landscape (0.56) = wide layout.",
			"4. SCALE: Apply scalingFormulas. For stories, bump ALL sizes one tier up: hero 96-120px, heading 56-72px, body 28-36px, card numbers 64-72px. Cards = 20%+ of canvas height. Section gaps = 64-96px.",
			"5. THINK IN CSS: Draft design as HTML/CSS first. Translate using cssToFigma map.",
			"6. HIERARCHY: Assign each element a level (primary/secondary/tertiary/ambient). Size, weight, color must reflect the hierarchy.",
			"7. CREATE: Multi-element = batch/bulk with named steps and interpolation. Single change = direct command.",
			"8. FRAME POSITIONING: New artboards go to the RIGHT of existing ones. Never overlap. Gap = 80-160px.",
			"9. BALANCE: ALL siblings MUST match — height (FIXED sizing), padding, spacing, radius, text sizes.",
			"10. EFFECTS: Shadows on cards/CTAs, gradients on hero backgrounds, blur for glass. See designDecisions.",
			"11. LAYERS: Name everything descriptively. Group in frames. Order back-to-front.",
			"12. EXPORT & VERIFY: export.image scale=2. Inspect the result. Does it look proportionate?",
		},
		"designPatterns": map[string]interface{}{
			"_overview":   "These patterns produce professional Figma designs. Follow them strictly.",
			"_decisionTree": "Use auto-layout for structured content (cards, lists, stacks, grids). Use coordinate math for decorative/absolute elements (overlays, background shapes, floating particles). Most designs use BOTH.",
			"coordinateSystem": map[string]interface{}{
				"origin":   "Top-left (0,0). X increases right, Y increases down.",
				"relative": "x/y are ALWAYS relative to the containing parent frame, not the page.",
				"rule":     "Inside auto-layout frames, x/y are IGNORED for auto-positioned children. Only absolute-positioned children use x/y.",
				"centering": map[string]interface{}{
					"horizontal": "x = (parentWidth - elementWidth) / 2",
					"vertical":   "y = (parentHeight - elementHeight) / 2",
					"both":       "x = (parentWidth - elementWidth) / 2, y = (parentHeight - elementHeight) / 2",
				},
				"edges": map[string]interface{}{
					"topLeft":     "x = margin, y = margin",
					"topRight":    "x = parentWidth - elementWidth - margin, y = margin",
					"bottomLeft":  "x = margin, y = parentHeight - elementHeight - margin",
					"bottomRight": "x = parentWidth - elementWidth - margin, y = parentHeight - elementHeight - margin",
				},
				"distribution": map[string]interface{}{
					"evenSpacing":  "gap = (parentWidth - N * itemWidth) / (N + 1); item[i].x = gap + i * (itemWidth + gap)",
					"spaceBetween": "gap = (parentWidth - N * itemWidth) / (N - 1); item[i].x = i * (itemWidth + gap)",
				},
				"grid": map[string]interface{}{
					"formula": "col = index % numColumns; row = floor(index / numColumns); x = marginLeft + col * (itemWidth + gapX); y = marginTop + row * (itemHeight + gapY)",
				},
				"snap8px": "Round all values to nearest multiple of 8: snapped = round(value / 8) * 8",
			},
			"grid": map[string]interface{}{
				"rule":    "Align ALL dimensions to an 8px grid. Every x, y, width, height, padding, and spacing value should be a multiple of 8.",
				"values":  "Common: 8, 16, 24, 32, 48, 64, 96, 128. Use 4 only for tight spacing.",
				"example": "width: 320 (not 310), padding: 24 (=3*8), spacing: 16 (=2*8), margin: 48 (=6*8).",
			},
			"framesNotShapes": map[string]interface{}{
				"rule":    "Use frames (node.create_frame) as containers, NOT rectangles. Frames support auto-layout, clipping, and nesting. A 'card' is a frame with auto-layout, not a rectangle with floating text.",
				"pattern": "1. Create parent frame. 2. Set auto-layout on it. 3. Add children (text, icons) INTO the frame using parentId param or layer.move_to_parent.",
			},
			"autoLayout": map[string]interface{}{
				"rule":      "Use auto-layout for ALL structural containers: cards, rows, columns, stacks, headers, footers. Auto-layout handles spacing and alignment — no manual x/y math needed for children.",
				"whenToUse": "Lists, stacks, card grids, buttons, navigation, header/content/footer, any group of siblings that need consistent spacing.",
				"centering": map[string]interface{}{
					"horizontal": "Set direction=VERTICAL, counterAxisAlign=CENTER. Children center horizontally.",
					"vertical":   "Set direction=VERTICAL, primaryAxisAlign=CENTER. Children center vertically.",
					"both":       "Set primaryAxisAlign=CENTER and counterAxisAlign=CENTER. Children center in both axes.",
				},
				"key_params": map[string]interface{}{
					"direction":         "HORIZONTAL or VERTICAL",
					"itemSpacing":       "Gap between children (use 8, 16, 24, 32)",
					"primaryAxisAlign":  "MIN (start), CENTER, MAX (end), SPACE_BETWEEN",
					"counterAxisAlign":  "MIN, CENTER, MAX",
					"padding":           "Internal padding (use 16, 24, 32, 48)",
					"primaryAxisSizing": "FIXED (exact size) or AUTO (hug contents)",
					"counterAxisSizing": "FIXED or AUTO",
				},
				"sizingModes": map[string]interface{}{
					"HUG":   "Frame shrinks to wrap its children tightly. Use for buttons, badges, tags.",
					"FILL":  "Child stretches to fill available parent space. Use for main content areas.",
					"FIXED": "Maintains exact dimensions. Use for root frames, cards with fixed size.",
				},
				"axisGuide": map[string]interface{}{
					"VERTICAL_frame":   "Primary axis = top-to-bottom. primaryAxisAlign controls vertical position. counterAxisAlign controls horizontal position.",
					"HORIZONTAL_frame": "Primary axis = left-to-right. primaryAxisAlign controls horizontal position. counterAxisAlign controls vertical position.",
				},
			},
			"absolutePositioning": map[string]interface{}{
				"rule":      "Use coordinate-based positioning for decorative elements, overlays, badges, and background shapes that float outside the auto-layout flow.",
				"whenToUse": "Floating decorations (rings, particles, lines), overlay badges, watermarks, background gradients, any element that overlaps or sits outside the content flow.",
				"howItWorks": []string{
					"Create the element and move it into the parent frame with layer.move_to_parent.",
					"The element will be positioned at the x/y coordinates RELATIVE to the parent frame.",
					"Use the coordinate math formulas to compute positions: centering, edges, corners.",
					"Decorative elements should be layered BEHIND content (use layer.send_to_back).",
				},
				"example_decorativeRing": `{"command":"shape.create_ellipse","params":{"x":10,"y":10,"width":160,"height":160,"name":"Deco Ring"}}, then set stroke, remove fill, move_to_parent into root frame`,
				"example_cornerBadge":    "x = parentWidth - badgeWidth - margin, y = margin. For 1080 frame with 32px badge and 16px margin: x = 1080 - 32 - 16 = 1032, y = 16",
			},
			"buildCard": map[string]interface{}{
				"description": "A card is a frame with auto-layout containing text children. NEVER create a rectangle + separate floating text.",
				"steps": []string{
					"1. node.create_frame {name, x, y, width, height} → save ID as 'card'",
					"2. paint.set_solid {nodeId: card, color, opacity} → background color",
					"3. node.set_corner_radius {nodeId: card, radius: 16}",
					"4. layout.set_auto_layout {nodeId: card, direction: VERTICAL, itemSpacing: 8, padding: 24}",
					"5. text.create {text: 'Title', parentId: card, fontSize: 18, fontWeight: 700}",
					"6. text.create {text: 'Description', parentId: card, fontSize: 14}",
				},
				"batch_example": `[
  {"name":"card","command":"node.create_frame","params":{"name":"Feature Card","x":0,"y":0,"width":320,"height":160}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.card.result.id}}","color":{"r":0.1,"g":0.1,"b":0.1}}},
  {"command":"node.set_corner_radius","params":{"nodeId":"${{steps.card.result.id}}","radius":16}},
  {"command":"layout.set_auto_layout","params":{"nodeId":"${{steps.card.result.id}}","direction":"VERTICAL","itemSpacing":8,"padding":24}},
  {"name":"title","command":"text.create","params":{"text":"Feature Title","parentId":"${{steps.card.result.id}}","fontSize":18}},
  {"command":"text.set_weight","params":{"nodeId":"${{steps.title.result.id}}","weight":700}},
  {"command":"text.set_color","params":{"nodeId":"${{steps.title.result.id}}","color":{"r":1,"g":1,"b":1}}},
  {"name":"desc","command":"text.create","params":{"text":"Description text","parentId":"${{steps.card.result.id}}","fontSize":14}},
  {"command":"text.set_color","params":{"nodeId":"${{steps.desc.result.id}}","color":{"r":0.6,"g":0.6,"b":0.6}}}
]`,
			},
			"layoutPatterns": map[string]interface{}{
				"description": "Common layout structures using nested auto-layout frames.",
				"centeredContentPage": map[string]interface{}{
					"description": "Full-page design with centered content stack. Use for social media posts, posters, etc.",
					"structure":   "Root frame (FIXED 1080x1080) → Content wrapper (auto-layout VERTICAL, CENTER/CENTER, padding 64) → children stack vertically and center automatically.",
					"batch_example": `[
  {"name":"root","command":"node.create_frame","params":{"name":"Post","x":0,"y":0,"width":1080,"height":1080}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.root.result.id}}","color":{"r":0.1,"g":0.1,"b":0.1}}},
  {"name":"content","command":"node.create_frame","params":{"name":"Content","parentId":"${{steps.root.result.id}}","width":1080,"height":1080}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.content.result.id}}","color":{"r":0,"g":0,"b":0,"a":0}}},
  {"command":"layout.set_auto_layout","params":{"nodeId":"${{steps.content.result.id}}","direction":"VERTICAL","itemSpacing":24,"padding":64,"primaryAxisAlign":"CENTER","counterAxisAlign":"CENTER"}}
]`,
				},
				"twoColumnGrid": map[string]interface{}{
					"description": "Two-column card row. Parent = HORIZONTAL auto-layout, children = card frames.",
					"batch_example": `[
  {"name":"row","command":"node.create_frame","params":{"name":"Card Row","width":680,"height":160}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.row.result.id}}","color":{"r":0,"g":0,"b":0,"a":0}}},
  {"command":"layout.set_auto_layout","params":{"nodeId":"${{steps.row.result.id}}","direction":"HORIZONTAL","itemSpacing":16}},
  {"name":"col1","command":"node.create_frame","params":{"name":"Col 1","parentId":"${{steps.row.result.id}}","width":328,"height":160}},
  {"name":"col2","command":"node.create_frame","params":{"name":"Col 2","parentId":"${{steps.row.result.id}}","width":328,"height":160}}
]`,
				},
				"headerContentFooter": map[string]interface{}{
					"description": "Page with header, scrollable content, and footer.",
					"structure":   "Root (VERTICAL, FIXED size) → Header (FIXED height) → Content (grows to fill) → Footer (FIXED height).",
				},
			},
			"typography": map[string]interface{}{
				"scale":      "Use consistent sizes from the scale: 12, 14, 16, 18, 20, 24, 32, 40, 48, 56, 64, 72, 80, 96.",
				"weights":    "400=Regular, 500=Medium, 600=SemiBold, 700=Bold, 800=ExtraBold.",
				"hierarchy":  "Hero: 64-96px/800wt. Heading: 32-48px/700wt. Subheading: 20-28px/500wt. Body: 16-20px/400wt. Caption: 12-14px/400wt.",
				"fontFamily": "Default to 'Inter' for UI. Set fontFamily explicitly for reliability.",
				"proportionality": map[string]interface{}{
					"rule":    "Font size MUST be proportional to canvas dimensions. Text that's too small will be unreadable on target devices. Use the minimum readable size for each format.",
					"formula": "minReadableSize = canvasWidth * 0.018 (roughly 2% of width). heroSize = canvasWidth * 0.06-0.08.",
					"socialMedia_1080px": map[string]interface{}{
						"hero":     "64-80px (6-7.5% of width). Use for main title.",
						"heading":  "36-48px. Use for section titles, feature numbers.",
						"body":     "20-28px. Use for descriptions, card titles. MINIMUM for readability.",
						"caption":  "16-20px. Smallest allowed text. Nothing below 16px on 1080px canvas.",
						"minSize":  "16px. Text below 16px is unreadable on mobile for 1080px social posts.",
					},
					"presentation_1920px": map[string]interface{}{
						"hero":    "80-120px.",
						"heading": "48-64px.",
						"body":    "28-36px.",
						"caption": "20-24px.",
						"minSize": "20px.",
					},
					"mobileUI_375px": map[string]interface{}{
						"hero":    "28-36px.",
						"heading": "20-24px.",
						"body":    "14-16px.",
						"caption": "11-12px.",
						"minSize": "11px.",
					},
					"tip": "After choosing font sizes, verify: is the smallest text at least 1.5% of canvas width? If not, increase it.",
				"aspectRatioRule": "Canvas aspect ratio changes BOTH proportions AND layout structure. Calculate heightRatio = height / width to decide layout strategy.",
				"layoutByAspectRatio": map[string]interface{}{
					"square_1to1": map[string]interface{}{
						"when":    "heightRatio ~1.0 (1080x1080, 1200x1200). Instagram posts, profile pictures.",
						"layout":  "Single centered content stack. Cards in 2-3 column rows. Compact vertical spacing.",
						"hero":    "64-80px. One hero, one subtitle, one card section, one footer.",
						"columns": "2-3 cards per row. Cards are wide and short (3:2 ratio).",
					},
					"portrait_9to16": map[string]interface{}{
						"when":    "heightRatio ~1.78 (1080x1920). Instagram/TikTok stories, Reels covers.",
						"layout":  "GENEROUS vertical spacing. Content fills middle 50-60%. Large top/bottom padding. More breathing room.",
						"hero":    "96-120px. Bigger impact. Add more vertical gap between sections (64-96px vs 32-48px for square).",
						"columns": "1-2 cards per row MAX. Cards are tall and wide (1:1 or taller). Stack vertically if content is heavy.",
						"critical": "Everything is bigger on stories. Bump ALL text one tier up. Cards should be 20% of height minimum.",
					},
					"landscape_16to9": map[string]interface{}{
						"when":    "heightRatio ~0.56 (1920x1080). Presentations, YouTube thumbnails, desktop banners.",
						"layout":  "Horizontal emphasis. Side-by-side layout. Content uses full width. Minimal top/bottom padding.",
						"hero":    "80-120px (based on height * 0.07-0.11). Wide text layouts.",
						"columns": "3-4 cards per row. Cards are shorter and wider (2:1 ratio). Use HORIZONTAL auto-layout for main content.",
					},
					"ultraTall_9to21": map[string]interface{}{
						"when":    "heightRatio >2.0 (1080x2400+). Some Android devices, vertical scrolling.",
						"layout":  "Stack everything vertically with VERY generous spacing. 1 card per row. Large section gaps.",
						"hero":    "120-160px. Massive impact needed to fill vertical space.",
						"columns": "1 column only. Full-width cards. Lots of vertical whitespace.",
					},
				},
				"scalingFormulas": map[string]interface{}{
					"_rule":       "Scale ALL dimensions relative to canvas size. Text uses width. Spacing uses BOTH dimensions. Tall canvases need bigger everything.",
					"heightRatio": "heightRatio = canvasHeight / canvasWidth. For squares: 1.0. For stories (1080x1920): 1.78. Multiply vertical dimensions by heightRatio.",
					"textSizing":  "Always based on WIDTH: hero = width * 0.07-0.09, heading = width * 0.04-0.05, body = width * 0.02-0.03, min = width * 0.018.",
					"storyTextBump": "For stories (height > width * 1.5): bump ALL font sizes up one tier. hero → 96-120px, heading → 56-72px, body → 28-36px, card numbers → 64-72px, card titles → 32-40px, card descriptions → 24-28px.",
					"padding":     "sidePadding = width * 0.06-0.08. topPadding = height * 0.15-0.25 (stories need more top space).",
					"cardHeight":  "cardHeight = height * 0.18-0.22 (346-422px on 1920px story). Cards should feel SUBSTANTIAL, not tiny.",
					"cardPadding": "cardPadding = width * 0.03-0.04 (32-43px on 1080px). For stories, use at least 36px.",
					"cardChildSpacing": "itemSpacing inside cards = 16-20px for stories (vs 8-12px for squares).",
					"sectionSpacing":   "Between major sections = height * 0.03-0.05 (58-96px on 1920px story). Between card rows = width * 0.02-0.04.",
					"contentFill":      "Content should fill 50-70% of canvas height. If content looks small, increase top/bottom padding to push content into the center.",
				},
				},
			},
			"colors": map[string]interface{}{
				"format":   "Use {r,g,b,a} objects where r/g/b are 0.0-1.0 floats, a is optional opacity.",
				"darkTheme": map[string]interface{}{
					"background":  "{r:0.1, g:0.1, b:0.1}",
					"surface":     "{r:1, g:1, b:1, a:0.04}",
					"border":      "{r:1, g:1, b:1, a:0.07}",
					"textPrimary": "{r:1, g:1, b:1}",
					"textMuted":   "{r:0.6, g:0.6, b:0.6}",
				},
				"tip": "For subtle glass/card effects, use low-opacity white fills with a 1px low-opacity white border.",
			},
			"shadows": map[string]interface{}{
				"subtle":  "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.15}, offsetX:0, offsetY:2, radius:8}",
				"medium":  "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.25}, offsetX:0, offsetY:4, radius:16}",
				"large":   "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.3}, offsetX:0, offsetY:8, radius:32}",
				"colored": "For accent glow: color:{r:1,g:0.84,b:0,a:0.15}, offsetX:0, offsetY:4, radius:24",
			},
			"exportQuality": map[string]interface{}{
				"rule":    "Always export at scale=2 (2x resolution) for crisp output. Default is 2x. Use scale=3 for print-quality exports.",
				"command": "export.image {nodeId, format:'PNG', scale:2}",
				"tip":     "After creating a design, export and verify it looks correct before reporting success.",
			},
			"balance": map[string]interface{}{
				"_rule": "Sibling elements MUST match in size, padding, and spacing. Unbalanced siblings look broken.",
				"evenCardHeights": map[string]interface{}{
					"problem": "Cards in a row have different heights because auto-layout HUGs content. Card with more text is taller.",
					"fix":     "Set primaryAxisSizing: FIXED on EVERY card in the row. All cards get the SAME fixed height.",
					"example": `Row of 3 cards, each 280x200:
  layout.set_sizing {nodeId: card1, primaryAxisSizing: FIXED}  // height stays 200
  layout.set_sizing {nodeId: card2, primaryAxisSizing: FIXED}  // height stays 200
  layout.set_sizing {nodeId: card3, primaryAxisSizing: FIXED}  // height stays 200

  WRONG: primaryAxisSizing: AUTO → cards HUG content → different heights`,
				},
				"consistentPadding": map[string]interface{}{
					"rule": "ALL sibling frames must share IDENTICAL padding values. Pick ONE padding for a group, apply to ALL.",
					"example": `3 cards in a row, all need padding 24:
  layout.set_auto_layout {nodeId: card1, padding: 24}
  layout.set_auto_layout {nodeId: card2, padding: 24}
  layout.set_auto_layout {nodeId: card3, padding: 24}

  WRONG: card1 padding:24, card2 padding:16, card3 padding:32`,
				},
				"consistentSpacing": map[string]interface{}{
					"rule": "ALL sibling frames must share IDENTICAL itemSpacing. Uneven gaps between children look broken.",
					"example": `3 cards, all need itemSpacing 12:
  layout.set_auto_layout {nodeId: card1, itemSpacing: 12}
  layout.set_auto_layout {nodeId: card2, itemSpacing: 12}
  layout.set_auto_layout {nodeId: card3, itemSpacing: 12}`,
				},
				"consistentCornerRadius": map[string]interface{}{
					"rule": "ALL sibling cards/containers share the SAME corner radius. Pick ONE radius, use everywhere.",
					"example": "3 cards → all cornerRadius: 16. WRONG: card1 radius:16, card2 radius:12, card3 radius:20.",
				},
				"siblingTextParity": map[string]interface{}{
					"rule":    "Text at the same hierarchy level across siblings must be the same fontSize, fontWeight, and color.",
					"example": "3 card titles → all fontSize:24, weight:700, color:white. 3 card descriptions → all fontSize:16, weight:400, color:muted.",
				},
				"widthDistribution": map[string]interface{}{
					"rule":    "Cards in a HORIZONTAL row should divide available width evenly. Use the same fixed width for each.",
					"formula": "cardWidth = (rowWidth - (numCards - 1) * gap) / numCards",
					"example": "Row 920px wide, 3 cards, gap 16: cardWidth = (920 - 2*16) / 3 = 296. All 3 cards get width 296.",
				},
				"checklist": []string{
					"All sibling cards: same width, same height (FIXED sizing)",
					"All sibling cards: same padding, same itemSpacing, same cornerRadius",
					"All sibling text at same level: same fontSize, weight, color",
					"Row gap matches design system (8, 16, 24, 32)",
					"Total row width = N * cardWidth + (N-1) * gap",
					"Card height should be at least 15-20% of canvas height (not tiny relative to canvas)",
					"Content area should use at least 60% of canvas height (not cramped at top)",
				},
			},
			"framePositioning": map[string]interface{}{
				"_rule": "When creating new root frames (artboards), NEVER overlap existing frames. Always check positions of existing frames first.",
				"workflow": []string{
					"1. BEFORE creating a new frame, call document.get_info or node.get_tree to see existing frames and their positions.",
					"2. Calculate the rightmost edge: maxX = max(frame.x + frame.width) for all existing root frames.",
					"3. Place new frame at x = maxX + gap (80-160px gap between artboards).",
					"4. For multi-frame sets (e.g., 3 story frames), space them: x = startX + i * (frameWidth + gap).",
				},
				"example": "Existing frame at x=0,w=1080. Next frame: x = 0 + 1080 + 80 = 1160. Third: x = 1160 + 1080 + 80 = 2320.",
			},
			"compositionTips": []string{
				"Create a root frame first (e.g., 1080x1080 for Instagram). All elements go inside it.",
				"Use auto-layout for STRUCTURED content (text stacks, card grids, navigation).",
				"Use coordinate math for DECORATIVE elements (background shapes, floating particles, accent lines).",
				"Build cards as frames with auto-layout, not as separate rectangles and text.",
				"Use layout.set_sizing with primaryAxisSizing=AUTO to let frames hug their content.",
				"Add padding to frames (24-48px) instead of offsetting child positions manually.",
				"For 2-column layouts, use a HORIZONTAL auto-layout frame with two child frames.",
				"Decorative elements go directly in the root frame with computed x/y coordinates.",
				"Use layer.send_to_back for background decorations so content appears on top.",
				"Name all elements descriptively — it helps with debugging and later editing.",
			},
		},
		"quickPrompts": []string{
			"For CREATING designs: build a batch JSON payload with named steps, interpolation, and auto-layout frames. Send with 'batch -f file.json'.",
			"For EDITING existing nodes: use single commands like paint.set_solid or node.resize.",
			"Use auto-layout frames as containers instead of absolute positioning.",
			"If one step fails, continue with remaining steps and return a structured summary.",
		},
		"tools": toolsOut,
	}
}

// SetupInstructions returns human-readable setup and installation instructions
// so the LLM can guide the user when the MCP is not yet connected to Figma.
func SetupInstructions() string {
	return `# AI Happy Design — Setup & Connection Guide

## Prerequisites
- Go 1.21+ (to build the binary)
- Node.js 18+ (to build the Figma plugin)
- Figma desktop app (not web — plugins need the desktop app)

## Step 1: Build the Binary
  cd <project-root>
  make build
This creates ./bin/ai-happy-design.

## Step 2: Build the Figma Plugin
  cd plugin && npm install && npm run build && cd ..

## Step 3: Load the Plugin in Figma
1. Open Figma desktop app
2. Go to Plugins > Development > Import plugin from manifest
3. Browse to <project-root>/plugin/manifest.json and select it
4. Run the plugin: Plugins > AI Happy Design
5. The plugin will show a UI panel — click Connect or it auto-connects

## Step 4: Configure MCP Client
Add this to your MCP client config (Claude Desktop, Cursor, etc.):

  {
    "mcpServers": {
      "ai-happy-design": {
        "command": "/absolute/path/to/bin/ai-happy-design",
        "args": ["mcp"]
      }
    }
  }

Replace /absolute/path/to/bin/ai-happy-design with the actual path.

## Step 5: Verify Connection
  ./bin/ai-happy-design tools --json
  ./bin/ai-happy-design command document.get_info

If document.get_info returns document data, the connection is working.

## Troubleshooting
- "Disconnected" in plugin: Make sure the relay is running (ai-happy-design mcp starts it automatically)
- Wrong port: The default is 3055. Edit Relay URL in the plugin UI if needed.
- Channel mismatch: Click Copy in the plugin UI to get the channel key, then pass it via --channel flag
- Plugin won't load: Rebuild with cd plugin && npm run build

## Once Connected
Call describe(action="catalog") to discover all available tools and start designing!`
}

// DesignGuide returns focused design guidance for LLM agents — design thinking,
// patterns, balance rules, and visual hierarchy without the full tool catalog.
// Call via describe(action="design_guide").
func DesignGuide() map[string]interface{} {
	catalog := LLMCatalog()
	return map[string]interface{}{
		"version":        catalog["version"],
		"designThinking": catalog["designThinking"],
		"designPatterns": catalog["designPatterns"],
		"playbook":       catalog["playbook"],
		"workflow":       catalog["workflow"],
		"hint":           "For full tool catalog with command examples, call describe(action='catalog').",
	}
}

func buildActionSpec(toolName, actionName, desc string) map[string]interface{} {
	description, paramsRaw := parseDescriptionAndParams(desc)
	paramNames := make([]string, 0, len(paramsRaw))
	sampleParams := map[string]interface{}{}
	for _, raw := range paramsRaw {
		norm := normalizeParamName(raw)
		if norm == "" {
			continue
		}
		paramNames = append(paramNames, norm)
		sampleParams[norm] = sampleValue(norm)
	}

	commandName := fmt.Sprintf("%s.%s", toolName, actionName)
	cliExample := buildCLIExample(commandName, sampleParams)
	mcpExample := buildMCPExample(toolName, actionName, sampleParams)
	batchStep := map[string]interface{}{
		"name":    fmt.Sprintf("%s_%s", toolName, actionName),
		"command": commandName,
		"params":  sampleParams,
	}

	// Curated overrides for higher-signal examples.
	if toolName == "paint" && actionName == "set_image_fill_from_url" {
		cliExample = "ai-happy-design command paint.set_image_fill_from_url -p '{\"nodeId\":\"1:2\",\"url\":\"https://picsum.photos/800/600\",\"scaleMode\":\"FILL\"}'"
		mcpExample = map[string]interface{}{
			"tool": "paint",
			"arguments": map[string]interface{}{
				"action":    "set_image_fill_from_url",
				"nodeId":    "1:2",
				"url":       "https://picsum.photos/800/600",
				"scaleMode": "FILL",
			},
		}
	}
	if toolName == "paint" && actionName == "set_image_fill" {
		cliExample = "ai-happy-design command paint.set_image_fill -p '{\"nodeId\":\"1:2\",\"imageData\":\"<base64-or-data-url>\",\"scaleMode\":\"FILL\"}'"
	}
	if toolName == "bulk" && actionName == "execute" {
		ops := `[{"name":"createCard","command":"shape.create_rectangle","params":{"x":40,"y":40,"width":220,"height":120}},{"name":"colorCard","command":"paint.set_solid","params":{"nodeId":"${{steps.createCard.result.id}}","color":"#2563EB"}}]`
		cliExample = fmt.Sprintf("ai-happy-design batch -o '%s' --retries 1 --retry-delay-ms 250", ops)
		mcpExample = map[string]interface{}{
			"tool": "bulk",
			"arguments": map[string]interface{}{
				"action":          "execute",
				"operations":      ops,
				"continueOnError": true,
				"retries":         1,
				"retryDelayMs":    250,
				"interpolate":     true,
			},
		}
		batchStep = map[string]interface{}{
			"name":    "bulk_execute",
			"command": "bulk.execute",
			"params": map[string]interface{}{
				"operations":      ops,
				"continueOnError": true,
				"retries":         1,
				"retryDelayMs":    250,
				"interpolate":     true,
			},
		}
	}

	spec := map[string]interface{}{
		"description": description,
		"params":      paramNames,
		"examples": map[string]interface{}{
			"cliCommand": cliExample,
			"mcpCall":    mcpExample,
			"batchStep":  batchStep,
		},
	}
	return spec
}

func buildCLIExample(commandName string, params map[string]interface{}) string {
	if len(params) == 0 {
		return fmt.Sprintf("ai-happy-design command %s", commandName)
	}

	body, err := json.Marshal(params)
	if err != nil {
		return fmt.Sprintf("ai-happy-design command %s", commandName)
	}
	return fmt.Sprintf("ai-happy-design command %s -p '%s'", commandName, string(body))
}

func buildMCPExample(toolName, actionName string, params map[string]interface{}) map[string]interface{} {
	args := map[string]interface{}{
		"action": actionName,
	}
	for key, val := range params {
		args[key] = val
	}

	return map[string]interface{}{
		"tool":      toolName,
		"arguments": args,
	}
}

func parseDescriptionAndParams(desc string) (string, []string) {
	parts := strings.SplitN(desc, "Params:", 2)
	if len(parts) < 2 {
		return strings.TrimSpace(desc), nil
	}

	description := strings.TrimSpace(parts[0])
	raw := strings.TrimSpace(parts[1])
	if raw == "" {
		return description, nil
	}
	if strings.EqualFold(raw, "No params") {
		return description, nil
	}

	tokens := strings.Split(raw, ",")
	params := make([]string, 0, len(tokens))
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		params = append(params, token)
	}
	return description, params
}

func normalizeParamName(raw string) string {
	name := strings.TrimSpace(raw)
	if name == "" {
		return ""
	}

	if strings.EqualFold(name, "No params") {
		return ""
	}

	// Keep the first token for entries like:
	// "padding (uniform) or paddingTop/Right/Bottom/Left"
	if idx := strings.IndexAny(name, " ("); idx > 0 {
		name = name[:idx]
	}

	name = strings.Trim(name, "*")

	// For slash variants, keep the first explicit key.
	if slash := strings.Index(name, "/"); slash > 0 {
		name = name[:slash]
	}

	return strings.TrimSpace(name)
}

func sampleValue(param string) interface{} {
	switch {
	case param == "x" || param == "y":
		return 40
	case strings.Contains(param, "width"):
		return 220
	case strings.Contains(param, "height"):
		return 120
	case strings.EqualFold(param, "scale"):
		return 1
	case strings.EqualFold(param, "opacity"):
		return 1.0
	case strings.Contains(strings.ToLower(param), "radius"):
		return 8
	case strings.Contains(strings.ToLower(param), "rotation"):
		return 0
	case strings.Contains(strings.ToLower(param), "spacing"):
		return 8
	case strings.Contains(strings.ToLower(param), "padding"):
		return 16
	case strings.Contains(strings.ToLower(param), "index"):
		return 0
	case strings.Contains(strings.ToLower(param), "visible") || strings.Contains(strings.ToLower(param), "locked"):
		return true
	case strings.HasSuffix(param, "Ids"):
		return []string{"1:2", "1:3"}
	case strings.HasSuffix(param, "Id"):
		return "1:2"
	case strings.EqualFold(param, "color"):
		return "#2563EB"
	case strings.EqualFold(param, "fontSize"):
		return 16
	case strings.EqualFold(param, "fontWeight"):
		return 500
	case strings.EqualFold(param, "fontFamily"):
		return "Inter"
	case strings.EqualFold(param, "fontStyle"):
		return "Regular"
	case strings.EqualFold(param, "name"):
		return "Example Node"
	case strings.EqualFold(param, "text"):
		return "Hello from AI Happy Design"
	case strings.EqualFold(param, "url"):
		return "https://picsum.photos/800/600"
	case strings.EqualFold(param, "imageData"):
		return "<base64-or-data-url>"
	case strings.EqualFold(param, "format"):
		return "PNG"
	case strings.EqualFold(param, "gradientType"):
		return "LINEAR"
	case strings.EqualFold(param, "stops"):
		return `[{"position":0,"color":"#2563EB"},{"position":1,"color":"#06B6D4"}]`
	case strings.EqualFold(param, "effects"):
		return `[{"type":"DROP_SHADOW","visible":true}]`
	case strings.EqualFold(param, "operations"):
		return `[{"name":"createRect","command":"shape.create_rectangle","params":{"x":40,"y":40,"width":220,"height":120}}]`
	case strings.EqualFold(param, "continueOnError"):
		return true
	case strings.EqualFold(param, "retries"):
		return 1
	case strings.EqualFold(param, "retryDelayMs"):
		return 250
	case strings.EqualFold(param, "interpolate"):
		return true
	default:
		return fmt.Sprintf("<%s>", param)
	}
}
