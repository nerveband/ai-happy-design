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
			"rule":      "CREATING = batch (one payload, many steps). EDITING = single command (precise, targeted).",
			"IMPORTANT": "ALWAYS call document.find_free_space BEFORE creating root frames. It returns exact x,y coordinates. NEVER hardcode or guess positions — this causes frames to overlap existing work.",
			"create": map[string]interface{}{
				"when": "Building new designs, layouts, multi-element compositions, or anything with 3+ elements.",
				"how":  "1) Call document.find_free_space to get x,y. 2) Build a JSON array of operations. First step = root frame at returned x,y. Subsequent steps can use compact refs like $frame or $last (long form ${{steps.name.result.id}} still works). 3) Prefer compact command aliases (frame/rect/ellipse/text/fill/stroke/parent/autolayout). 4) Send: 'ai-happy-design batch '[...]'' or bulk.execute via MCP.",
				"why":  "One WebSocket connection, no process overhead per step, automatic ID chaining. 150 steps in ~6 seconds vs ~8 minutes with individual commands.",
			},
			"edit": map[string]interface{}{
				"when": "Changing a color, moving a node, resizing, renaming, or any single-property change on an existing node.",
				"how":  "Use a single command: ai-happy-design command paint.set_solid '{\"nodeId\":\"1:2\",\"color\":\"#FF0000\"}'",
				"why":  "Fast, precise, no batch overhead needed for one operation.",
			},
			"multiBatch": map[string]interface{}{
				"when": "Creating multiple independent designs (e.g., 6 carousel slides, a set of social posts). Each design is its own batch file.",
				"how": []string{
					"Write each design as a separate .json batch file.",
					"Run all at once: ai-happy-design batch slide1.json slide2.json slide3.json",
					"Or use a directory: ai-happy-design batch ./slides/",
					"Or a glob: ai-happy-design batch slides/*.json",
					"Add --parallel for concurrent execution (up to 4 at a time).",
				},
				"parallel": "ai-happy-design batch slide1.json slide2.json --parallel — each file gets its own WebSocket connection and auto-placement.",
				"why":      "No need for Python scripts or shell loops. Just write JSON files and batch them.",
			},
			"connectionDiagnostics": map[string]interface{}{
				"_overview": "The CLI auto-checks plugin connectivity before sending commands. If no plugin is connected, it shows a diagnostic message and waits up to 30s for the plugin to connect. This prevents the 300s silent hang.",
				"behavior":  "If relay is running but no plugin is connected, you'll see: 'Relay running but no Figma plugin connected on channel X. Waiting up to 30s...' — then either the plugin connects and the command proceeds, or it fails with a clear error.",
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
				"_hint":                            "When you think 'I would use this CSS', translate it to the corresponding Figma command.",
				"display: flex":                    "layout.set_auto_layout {direction, itemSpacing, padding}",
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
						"how":     "Use display/hero tier from compute_tokens, boldest weight (800), accent color, optional shadow glow.",
						"example": "Use text.display or text.hero from compute_tokens (e.g. 200px or 152px on 1080w). Weight 800, accent color.",
					},
					"secondary": map[string]interface{}{
						"what":    "Supporting elements. Subtitles, section headers, card titles, secondary CTAs.",
						"how":     "Use heading/subheading tier from compute_tokens, semi-bold (600-700), white or slightly muted color.",
						"example": "Use text.heading or text.subheading from compute_tokens (e.g. 84px or 64px on 1080w). Weight 600-700.",
					},
					"tertiary": map[string]interface{}{
						"what":    "Details. Descriptions, body text, metadata, helper text.",
						"how":     "Use body/caption tier from compute_tokens, regular weight (400-500), muted color (gray/low opacity).",
						"example": "Use text.body or text.caption from compute_tokens (e.g. 48px or 36px on 1080w). Weight 400-500, color:{r:0.6,g:0.6,b:0.6}.",
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
					"recipe_subtle":   "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.1}, offsetX:0, offsetY:2, radius:8}",
					"recipe_card":     "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.2}, offsetX:0, offsetY:4, radius:16}",
					"recipe_elevated": "effect.add_shadow {shadowType:DROP_SHADOW, color:{r:0,g:0,b:0,a:0.3}, offsetX:0, offsetY:8, radius:32}",
					"recipe_glow":     "effect.add_shadow {shadowType:DROP_SHADOW, color:{accent_r,accent_g,accent_b,a:0.2}, offsetX:0, offsetY:4, radius:24}",
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
					"recipe_bgSweep": "paint.set_gradient {type:LINEAR, stops:[{position:0,color:dark},{position:1,color:slightlyLighter}]}",
					"recipe_accent":  "paint.set_gradient {type:LINEAR, stops:[{position:0,color:accent},{position:1,color:accentVariant}]}",
					"recipe_fade":    "paint.set_gradient {type:LINEAR, stops:[{position:0,color:{r:1,g:1,b:1,a:0.1}},{position:1,color:{r:1,g:1,b:1,a:0}}]}",
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
					"recipe_subtle":   "paint.set_stroke {color:{r:1,g:1,b:1,a:0.07}, width:1}",
					"recipe_accent":   "paint.set_stroke {color:{accent,a:0.12}, width:1}",
					"recipe_decoRing": "paint.set_stroke {color:{accent,a:0.1}, width:2} + paint.remove_fill (stroke-only)",
				},
			},
			"layerOrganization": map[string]interface{}{
				"naming": map[string]interface{}{
					"rule":    "Name every element by its content or role. Never leave default names like 'Frame 47' or 'Rectangle 12'.",
					"pattern": "[Type] - [Purpose]",
					"examples": map[string]interface{}{
						"frames": "Hero Section, Card Grid, Card - Feature Name, Footer, Navigation Bar",
						"text":   "Hero Title, Card Title, Card Description, Footer Link, Badge Label",
						"shapes": "Deco Ring, Background Gradient, Divider, Dot Indicator, Accent Line",
						"images": "Avatar - User, Hero Background, Logo, Product Thumbnail",
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
			"2. COMPUTE TOKENS: Call design.compute_tokens with your canvas width and height. This returns ALL sizing values (font sizes, spacing, padding, layout type, card widths) as concrete pixel values. Use these EXACT values — do NOT calculate sizes yourself.",
			"3. FIND SPACE: Call document.find_free_space with your desired width/height to get exact x,y coordinates for new frames. NEVER guess frame positions.",
			"4. THINK IN CSS: Draft the design as HTML/CSS in your head. The compute_tokens response includes a CSS analogy. Translate to Figma using the cssToFigma map.",
			"5. HIERARCHY: Assign each element a level. Use the modular type scale from compute_tokens: display=hero stats/amounts, hero=campaign headlines, title=slide headlines, heading=section titles, subheading=labels, body=copy, caption=footnotes. Use button tokens for CTA sizing.",
			"6. CREATE: Multi-element = batch/bulk with named steps and interpolation. Single change = direct command.",
			"7. FRAME POSITIONING: Call document.find_free_space to get exact coordinates. NEVER hardcode or guess x/y for root frames.",
			"8. BALANCE: ALL siblings MUST match — padding, spacing, radius, text sizes from compute_tokens.",
			"9. TEXT WIDTH: After creating text inside auto-layout, resize text to contentWidth from compute_tokens.",
			"10. EFFECTS: Shadows on cards/CTAs, gradients on hero backgrounds, blur for glass. See designDecisions.",
			"11. LAYERS: Name everything descriptively. Group in frames. Order back-to-front.",
			"12. EXPORT & VERIFY: export.image scale=2. Inspect the result. Does it look proportionate?",
		},
		"designPatterns": map[string]interface{}{
			"_overview":     "These patterns produce professional Figma designs. Follow them strictly.",
			"_decisionTree": "DEFAULT: Use auto-layout (flexbox). Think in rows and columns. Figma auto-layout IS CSS flexbox — you already know it. Create frames with layoutMode, itemSpacing, padding, and alignment in a SINGLE create_frame call. For decorative overlays inside auto-layout frames, use layoutPositioning:ABSOLUTE on the child to exempt it from the flow. After creating, run layout.check_overlaps to verify no elements overlap.",
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
				"rule":      "Use auto-layout for ALL layout containers. Auto-layout IS CSS flexbox. Figma's layout engine handles positioning — no manual x/y math needed.",
				"whenToUse": "Everything: page layouts, cards, rows, columns, stacks, headers, footers, buttons, badges. The ONLY exception is decorative overlays.",
				"flexboxMapping": map[string]interface{}{
					"_hint":                          "You already know CSS flexbox. Use this mapping:",
					"flex-direction: column":         "layoutMode: VERTICAL",
					"flex-direction: row":            "layoutMode: HORIZONTAL",
					"gap":                            "itemSpacing",
					"padding":                        "padding (uniform) or paddingTop/Right/Bottom/Left",
					"justify-content: flex-start":    "primaryAxisAlignItems: MIN",
					"justify-content: center":        "primaryAxisAlignItems: CENTER",
					"justify-content: flex-end":      "primaryAxisAlignItems: MAX",
					"justify-content: space-between": "primaryAxisAlignItems: SPACE_BETWEEN",
					"align-items: flex-start":        "counterAxisAlignItems: MIN",
					"align-items: center":            "counterAxisAlignItems: CENTER",
					"align-items: flex-end":          "counterAxisAlignItems: MAX",
					"width/height on main axis":      "primaryAxisSizingMode: FIXED or AUTO",
					"cross-axis size":                "counterAxisSizingMode: FIXED or AUTO",
					"flex-grow: 1":                   "child layoutGrow: 1 (only 0 or 1, no fractions)",
					"align-self: stretch":            "child layoutAlign: STRETCH",
					"width: fit-content":             "child layoutSizingHorizontal: HUG",
					"width: 100%":                    "child layoutSizingHorizontal: FILL",
					"explicit width":                 "child layoutSizingHorizontal: FIXED",
					"flex-wrap: wrap":                "layoutWrap: WRAP",
					"position: absolute":             "child layoutPositioning: ABSOLUTE (for decorative elements only)",
				},
				"textInAutoLayout": map[string]interface{}{
					"_critical": "Text nodes inside auto-layout need special handling to avoid height collapse.",
					"problem":   "textAutoResize=HEIGHT alone does NOT work inside auto-layout. Text collapses because auto-layout controls sizing.",
					"fix":       "Our plugin AUTO-SETS layoutSizingVertical=HUG when text is in auto-layout with textAutoResize=HEIGHT. You don't need to do anything extra.",
					"sequence":  "1. Create text with width param (enables HEIGHT resize). 2. Plugin auto-sets layoutSizingVertical=HUG. 3. Optionally set layoutSizingHorizontal=FILL for full-width text.",
				},
				"oneCommandCreate": map[string]interface{}{
					"_hint":   "Create auto-layout frames in ONE command — no separate layout.set_auto_layout call needed:",
					"example": `node.create_frame {name:"VStack", width:500, layoutMode:"VERTICAL", itemSpacing:24, padding:32, primaryAxisAlign:"CENTER", counterAxisAlign:"CENTER", primaryAxisSizing:"AUTO"}`,
				},
				"spacers": map[string]interface{}{
					"fixedSpacer": "Create empty frame with fixed height (e.g., 48px) as a child. Acts like a CSS margin.",
					"flexSpacer":  "Create empty frame with layoutGrow:1. Pushes siblings apart like flex-grow spacer.",
				},
				"badge": map[string]interface{}{
					"recipe": `node.create_frame {name:"Badge", layoutMode:"HORIZONTAL", padding:8, paddingLeft:16, paddingRight:16, primaryAxisSizing:"AUTO", counterAxisSizing:"AUTO", color:"#FF6B00"} + node.set_corner_radius {radius:100}`,
				},
				"decorativeAbsolute": map[string]interface{}{
					"_hint":   "For decorative elements (circles, stripes, gradients) inside auto-layout frames, set layoutPositioning:ABSOLUTE on the child. This takes it out of the flow — like CSS position:absolute.",
					"example": "shape.create_ellipse {parentId:frameId, x:0, y:0, width:200, height:200, layoutPositioning:ABSOLUTE}",
				},
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
				"gradientTransforms": map[string]interface{}{
					"leftToRight": "[[1,0,0],[0,1,0]]",
					"topToBottom": "[[0,1,0],[-1,0,1]]",
				},
			},
			"manualPositioning": map[string]interface{}{
				"_overview":      "For non-auto-layout frames, children are positioned by x/y coordinates relative to the parent. No special property is needed — this is Figma's default behavior. Do NOT set layoutPositioning:ABSOLUTE on children of non-auto-layout frames (it will error).",
				"rule":           "Non-auto-layout frame + children with x/y = manual positioning. Auto-layout frame + layoutPositioning:ABSOLUTE = exempt child from flow (decorative overlays only).",
				"WARNING":        "layoutPositioning:ABSOLUTE is ONLY valid inside auto-layout parents. Setting it on a child of a frame with layoutMode:NONE causes an error. If the parent has no auto-layout, just use x/y — no extra property needed.",
				"colorShortcut":  "node.create_frame, shape.create_rectangle, and text.create all accept a 'color' param (hex string like '#0F0F23' or rgb object). The param name is 'color', NOT 'fillColor' or 'backgroundColor'. This is a shortcut that avoids a separate paint.set_solid call.",
				"howItWorks": []string{
					"1. Create root frame (NO auto-layout): node.create_frame {name, x, y, width, height, color: '#0F0F23'}",
					"2. Add children with parentId and x/y: text.create {text, parentId, x, y, width, fontSize, color: '#FFFFFF'}",
					"3. Children are positioned by x/y automatically — NO layoutPositioning needed.",
					"4. For badges/buttons that need centered text: use auto-layout on the badge frame itself.",
					"5. x/y are RELATIVE to the parent frame, not the page.",
					"6. Text width: use the 'width' parameter on text.create for text wrapping control.",
				},
				"coordinatePlanning": []string{
					"Plan your layout TOP-DOWN. Write down the y-coordinate of each major element.",
					"Account for text height: height ≈ fontSize * (lineHeight/100) * numLines.",
					"Leave 48-96px gaps between major sections. Use compute_tokens for spacing values.",
					"Cards: use consistent heights for siblings (e.g., all 200px). Place with exact x/y.",
					"Horizontal card rows: cardWidth = (contentWidth - (N-1) * gap) / N, then x[i] = sidePadding + i * (cardWidth + gap).",
				},
				"hybridExample": "Root frame (manual x/y) → Badge (auto-layout for centering) → Hero text (x/y, width param) → Subtitle (x/y, width param) → Cards (x/y) → CTA button (auto-layout for centering)",
			},
			"buildCard": map[string]interface{}{
				"description": "A card is a frame with text children positioned using x/y. NEVER create a rectangle + separate floating text.",
				"steps": []string{
					"1. node.create_frame {name, parentId, x, y, width, height} → save ID as 'card'",
					"2. paint.set_solid {nodeId: card, color:{r:1,g:1,b:1}, opacity:0.04} → glass fill",
					"3. node.set_corner_radius {nodeId: card, radius: 16}",
					"4. paint.set_stroke {nodeId: card, color:{r:1,g:1,b:1}, opacity:0.07, strokeWeight:1}",
					"5. text.create {text: 'Title', parentId: card, x:padding, y:padding, fontSize: tokens.text.heading, fontWeight: 700}",
					"6. text.create {text: 'Description', parentId: card, x:padding, y:titleBottom+8, fontSize: tokens.text.body, width: cardWidth-2*padding}",
				},
				"batch_example": `[
  {"name":"card","command":"node.create_frame","params":{"name":"Feature Card","parentId":"${{steps.root.result.id}}","x":72,"y":440,"width":288,"height":200}},
  {"command":"paint.set_solid","params":{"nodeId":"${{steps.card.result.id}}","color":{"r":1,"g":1,"b":1},"opacity":0.04}},
  {"command":"node.set_corner_radius","params":{"nodeId":"${{steps.card.result.id}}","radius":16}},
  {"command":"paint.set_stroke","params":{"nodeId":"${{steps.card.result.id}}","color":{"r":1,"g":1,"b":1},"opacity":0.07,"strokeWeight":1}},
  {"name":"title","command":"text.create","params":{"text":"10x","parentId":"${{steps.card.result.id}}","x":32,"y":36,"fontSize":48,"fontWeight":700,"fontFamily":"Space Grotesk"}},
  {"command":"text.set_color","params":{"nodeId":"${{steps.title.result.id}}","color":{"r":0.6,"g":0.4,"b":1}}},
  {"name":"desc","command":"text.create","params":{"text":"Faster Design","parentId":"${{steps.card.result.id}}","x":32,"y":128,"fontSize":28,"fontWeight":500,"fontFamily":"DM Sans"}},
  {"command":"text.set_color","params":{"nodeId":"${{steps.desc.result.id}}","color":{"r":0.6,"g":0.6,"b":0.68}}}
]`,
			},
			"fillRules": map[string]interface{}{
				"_critical": "ONLY visual surfaces should have fills. Structural/layout frames (wrappers, groups, auto-layout containers) should have NO fill.",
				"visualSurfaces": []string{
					"Slide backgrounds (root frames) — use paint.set_solid AFTER creating the frame",
					"Cards and elevated surfaces — solid or glass fill",
					"Buttons — accent color fill",
					"Badges — colored background for contrast",
				},
				"structuralFrames": []string{
					"Content wrappers — NO fill (remove default white fill with paint.remove_fill {index:0})",
					"Auto-layout groups — NO fill",
					"Header/Content/Footer groups — NO fill",
					"Row/column containers — NO fill",
				},
				"frameFillParam": "node.create_frame accepts a 'color' param that sets the fill on creation. Use this for VISUAL frames (slide bg, cards). For structural frames, either don't pass 'color' and call paint.remove_fill, or simply leave the default.",
				"removeFill":     "paint.remove_fill {nodeId, index:0} — removes the fill at index 0. Use on structural frames to clear default white fill. Command is 'paint.remove_fill' (singular), NOT 'paint.remove_fills'.",
			},
			"layoutPatterns": map[string]interface{}{
				"description": "Common layout structures using nested auto-layout frames.",
				"centeredContentPage": map[string]interface{}{
					"description": "Full-page design with centered content stack. Use for social media posts, posters, etc.",
					"structure":   "Root frame (FIXED 1080x1080, with bg fill) → Content wrapper (auto-layout VERTICAL, CENTER/CENTER, NO fill) → children stack vertically and center automatically.",
					"batch_example": `[
  {"name":"root","command":"node.create_frame","params":{"name":"Post","x":0,"y":0,"width":1080,"height":1080,"color":"#1A1A1A"}},
  {"name":"content","command":"node.create_frame","params":{"name":"Content","parentId":"${{steps.root.result.id}}","width":1080,"height":1080}},
  {"command":"paint.remove_fill","params":{"nodeId":"${{steps.content.result.id}}","index":0}},
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
				"weights":     "400=Regular, 500=Medium, 600=SemiBold, 700=Bold, 800=ExtraBold.",
				"fontFamily":  "Default to 'Inter' for clean UI. Set fontFamily explicitly for reliability.",
				"googleFonts": "Figma includes ALL Google Fonts by default. You can use any Google Font family name (e.g., 'Poppins', 'Space Grotesk', 'Playfair Display', 'DM Sans', 'Outfit', 'Sora', 'Manrope'). Great for adding personality to designs.",
				"lineHeight":  "CRITICAL: You must specify lineHeightUnit:'PERCENT' when calling text.set_line_height. Example: {lineHeight: 130, lineHeightUnit: 'PERCENT'} for 130%. Without lineHeightUnit, the value is interpreted as PIXELS (e.g., 130px), causing massive text overflow. Hero: 110%, Body: 140%, Caption: 130%.",
				"fontPairings": map[string]interface{}{
					"modern":  "Headings: 'Space Grotesk' or 'Outfit'. Body: 'Inter' or 'DM Sans'.",
					"elegant": "Headings: 'Playfair Display'. Body: 'Source Sans Pro' or 'Lato'.",
					"bold":    "Headings: 'Sora' or 'Manrope'. Body: 'Inter'.",
					"tech":    "Headings: 'JetBrains Mono'. Body: 'Inter' or 'IBM Plex Sans'.",
				},
			},
			"sizingSystem": map[string]interface{}{
				"_philosophy": "Think of each canvas like an HTML viewport. You already know CSS — use that knowledge. A 1080x1920 story is like a mobile viewport in portrait. A 1920x1080 canvas is like a desktop viewport. Size text the same way you would with CSS rem/vw units, but output concrete pixel values.",
				"_method":     "BEST: Call design.compute_tokens {width, height} for exact values. This table is a quick reference using the modular scale (base = W * 0.044, ratio = 1.333). Text on 4px grid, spacing on 8px grid.",
				"lookupTable": map[string]interface{}{
					"_instruction": "PREFER calling design.compute_tokens. This table is a quick reference. Scale: caption(-1) → body(0) → subheading(+1) → heading(+2) → title(+3) → hero(+4) → display(+5).",
					"W_1080": map[string]interface{}{
						"_use":    "Instagram post (1080x1080), Instagram story (1080x1920), Instagram reel cover, TikTok, most social media",
						"display": 200, "hero": 152, "title": 112, "heading": 84, "subheading": 64, "body": 48, "caption": 36, "numbers": 200, "cta": 48,
						"sidePadding": 72, "contentWidth": 936, "cardPadding": 40, "framePadding": 64,
						"itemSpacing": 16, "cardGap": 24, "cornerRadius": 16,
					},
					"W_1200": map[string]interface{}{
						"_use":    "Facebook post (1200x630), LinkedIn post (1200x627), larger social cards",
						"display": 224, "hero": 168, "title": 124, "heading": 92, "subheading": 72, "body": 52, "caption": 40, "numbers": 224, "cta": 52,
						"sidePadding": 80, "contentWidth": 1040, "cardPadding": 40, "framePadding": 72,
						"itemSpacing": 16, "cardGap": 24, "cornerRadius": 16,
					},
					"W_1920": map[string]interface{}{
						"_use":    "Presentation (1920x1080), desktop banner",
						"display": 356, "hero": 268, "title": 200, "heading": 152, "subheading": 112, "body": 84, "caption": 64, "numbers": 356, "cta": 84,
						"sidePadding": 128, "contentWidth": 1664, "cardPadding": 72, "framePadding": 112,
						"itemSpacing": 24, "cardGap": 40, "cornerRadius": 24,
					},
					"W_1280": map[string]interface{}{
						"_use":    "YouTube thumbnail (1280x720), medium presentations",
						"display": 236, "hero": 180, "title": 132, "heading": 100, "subheading": 76, "body": 56, "caption": 44, "numbers": 236, "cta": 56,
						"sidePadding": 88, "contentWidth": 1104, "cardPadding": 48, "framePadding": 80,
						"itemSpacing": 24, "cardGap": 32, "cornerRadius": 24,
					},
					"W_800": map[string]interface{}{
						"_use":    "Email header, blog hero image, medium-sized graphics",
						"display": 148, "hero": 112, "title": 84, "heading": 64, "subheading": 48, "body": 36, "caption": 28, "numbers": 148, "cta": 36,
						"sidePadding": 56, "contentWidth": 688, "cardPadding": 24, "framePadding": 48,
						"itemSpacing": 16, "cardGap": 16, "cornerRadius": 16,
					},
					"W_custom": map[string]interface{}{
						"_use":  "Call design.compute_tokens {width, height} for exact values. Or: base = W * 0.044, scale by 1.333 per step. Round text to 4px, spacing to 8px.",
						"scale": "caption = base/1.333, body = base, subheading = base×1.333, heading = base×1.333², title = base×1.333³, hero = base×1.333⁴, display = base×1.333⁵",
					},
				},
				"layoutByAspectRatio": map[string]interface{}{
					"_rule": "After looking up sizes from the table, determine LAYOUT STRUCTURE from the aspect ratio. heightRatio = H / W.",
					"square": map[string]interface{}{
						"when":    "heightRatio 0.9 - 1.1 (e.g., 1080x1080, 1200x1200)",
						"columns": "2-3 cards per row",
						"layout":  "Think: a centered landing page section. Vertical stack with hero, subtitle, card grid, CTA. Like a single viewport with padding: 8% top.",
						"css":     "body { display:flex; flex-direction:column; align-items:center; padding:8% 6.5%; gap:3.5%; }",
					},
					"portrait": map[string]interface{}{
						"when":    "heightRatio 1.3 - 2.0 (e.g., 1080x1920 story, 1080x1350 post)",
						"columns": "1-2 cards per row MAX",
						"layout":  "Think: a full-screen mobile app. Generous vertical spacing. Content centered in middle 60%. Big breathing room top/bottom. Like a centered hero section on mobile.",
						"css":     "body { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:12% 6.5%; gap:4%; }",
						"tip":     "Portrait has LOTS of vertical space. Use it. Section gaps should be bigger. Content should feel centered and balanced, never cramped at the top.",
					},
					"landscape": map[string]interface{}{
						"when":    "heightRatio 0.4 - 0.7 (e.g., 1920x1080, 1280x720)",
						"columns": "3-4 cards per row",
						"layout":  "Think: a desktop hero banner. Side-by-side layouts. Content uses full width. Compact vertically. Like a wide navbar + hero section.",
						"css":     "body { display:flex; flex-direction:row; align-items:center; padding:6% 6.5%; gap:2%; }",
						"tip":     "Height is the constraint. Keep content compact vertically. Use a two-column absolute layout: left side = text content, right side = stat cards.",
					},
				},
				"verticalSpacing": map[string]interface{}{
					"_rule":    "Section gaps scale with canvas HEIGHT. Use these values based on H.",
					"H_1080":   map[string]interface{}{"sectionGap": 40, "topPadding": 88},
					"H_1920":   map[string]interface{}{"sectionGap": 64, "topPadding": 232},
					"H_1350":   map[string]interface{}{"sectionGap": 48, "topPadding": 160},
					"H_720":    map[string]interface{}{"sectionGap": 24, "topPadding": 40},
					"H_custom": "sectionGap = round8(H * 0.035). topPadding for portrait = round8(H * 0.12), for square = round8(H * 0.08), for landscape = round8(H * 0.06).",
				},
			},
			"colors": map[string]interface{}{
				"format": "Use {r,g,b,a} objects where r/g/b are 0.0-1.0 floats, a is optional opacity.",
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
			"printDesign": map[string]interface{}{
				"_overview": "For print designs (flyers, posters, business cards), pass dpi parameter to design.compute_tokens.",
				"dpiGuide": map[string]interface{}{
					"screen":     "72 dpi (default). For social media, web, presentations.",
					"print":      "300 dpi. For professional print — flyers, posters, brochures, business cards.",
					"largePrint": "150 dpi. For large format — banners, billboards (viewed from distance).",
				},
				"usage": "design.compute_tokens {width:2550, height:3300, dpi:300} → body ≈ 50px (12pt at 300dpi), proper print sizing",
				"physicalSizes": map[string]interface{}{
					"letter":       "2550×3300 @ 300dpi = 8.5×11 inches",
					"A4":           "2480×3508 @ 300dpi = 8.27×11.69 inches",
					"businessCard": "1050×600 @ 300dpi = 3.5×2 inches",
					"poster":       "5400×7200 @ 300dpi = 18×24 inches",
				},
				"fontGuideline": "Print body text: 10-14pt. Headers scale up from there. The dpi parameter handles all conversions automatically.",
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
					"rule":    "ALL sibling cards/containers share the SAME corner radius. Pick ONE radius, use everywhere.",
					"example": "3 cards → all cornerRadius: 16. WRONG: card1 radius:16, card2 radius:12, card3 radius:20.",
				},
				"siblingTextParity": map[string]interface{}{
					"rule":    "Text at the same hierarchy level across siblings must be the same fontSize, fontWeight, and color.",
					"example": "3 card titles → all fontSize:heading (e.g. 48), weight:700, color:white. 3 card descriptions → all fontSize:body (e.g. 28), weight:400, color:muted.",
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
				"_rule": "NEVER guess frame positions. ALWAYS call document.find_free_space BEFORE creating any root frame.",
				"workflow": []string{
					"1. Call document.find_free_space with width, height, and gap. It returns exact x,y coordinates and a list of all existing frames.",
					"2. Use the returned x,y as the position for your new frame. These are 8px-grid-snapped and guaranteed non-overlapping.",
					"3. For multi-frame sets, call find_free_space once, then offset additional frames: x = returned_x + i * (frameWidth + gap).",
				},
				"example": "document.find_free_space {width: 1080, height: 1920, gap: 100} → {x: 1280, y: 0}. Use x=1280, y=0 for your frame.",
			},
			"textRules": map[string]interface{}{
				"_critical": "Use text.create with a 'width' parameter to control text wrapping. Always set lineHeightUnit:'PERCENT'.",
				"textWidth": []string{
					"BEST: Set width on text.create: text.create {text: '...', x: 72, y: 100, width: 936, fontSize: 24}",
					"The width parameter on text.create sets the text box width directly. Text auto-wraps within it.",
					"This works for BOTH absolute-positioned and auto-layout text.",
				},
				"autoLayoutText": []string{
					"For text inside auto-layout, you can ALSO use: layout.set_sizing {nodeId: textId, horizontal: 'FILL'}",
					"This makes text fill the parent width (like CSS width:100% in flexbox).",
				},
				"lineHeight": []string{
					"CRITICAL: Always pass lineHeightUnit:'PERCENT' when setting line height.",
					"Example: text.set_line_height {nodeId: textId, lineHeight: 140, lineHeightUnit: 'PERCENT'}",
					"WITHOUT lineHeightUnit, the value is treated as PIXELS (e.g., 140px line height), causing massive text overflow.",
					"Good values: hero=110%, heading=120%, body=140%, caption=130%.",
				},
				"newlines": []string{
					"NEVER use \\n to split text into separate visual lines. Figma treats \\n as a paragraph break with extra spacing.",
					"Instead, create separate text nodes for each distinct text block.",
					"For multi-line paragraphs, use a single text node — it auto-wraps within its width.",
				},
			},
			"cardRules": map[string]interface{}{
				"_critical": "Cards are frames with text children positioned inside using x/y. Use consistent heights for sibling cards.",
				"construction": []string{
					"Call design.compute_tokens first to get cardPadding, cornerRadius, and card widths.",
					"Create card frame: node.create_frame {name, parentId, x, y, width, height}. Use widths from compute_tokens.",
					"Set a fixed height (e.g., 176-200px) — all sibling cards MUST have the same height.",
					"Set cornerRadius from compute_tokens. Add fill {r:1,g:1,b:1,a:0.04} and stroke {r:1,g:1,b:1,a:0.07}.",
					"Add text children with explicit x/y RELATIVE to the card: text.create {text, parentId:cardId, x:padding, y:padding, fontSize:tokens.text.heading}.",
					"Use font sizes from compute_tokens: text.heading for titles, text.body for descriptions, text.caption for metadata.",
				},
			},
			"compositionTips": []string{
				"FIRST: Call design.compute_tokens {width, height} to get all sizing values. Use them throughout.",
				"Create a root frame with auto-layout: node.create_frame {name, width, height, layoutMode:VERTICAL, padding, itemSpacing, primaryAxisAlign:CENTER, counterAxisAlign:CENTER}",
				"DEFAULT APPROACH: Use auto-layout (flexbox) for ALL layout. Think in rows and columns.",
				"Text in auto-layout: set width param on text.create. Plugin auto-sets layoutSizingVertical:HUG.",
				"CRITICAL: Always set lineHeightUnit:'PERCENT' when calling text.set_line_height. Default is PIXELS which breaks layouts.",
				"layoutPositioning:ABSOLUTE is ONLY for decorative overlays inside AUTO-LAYOUT frames. Never use it on children of non-auto-layout frames — they already position by x/y.",
				"Use padding and spacing values from compute_tokens — never guess pixel values.",
				"After creating layout, call layout.check_overlaps {nodeId} to verify no elements overlap.",
				"Name all elements descriptively. Layer background decorations behind content with layer.send_to_back.",
				"For batch export: export.batch {nodeIds:'id1,id2', format:'PNG', scale:2}",
			},
		},
		"advancedFeatures": map[string]interface{}{
			"modify": map[string]interface{}{
				"_overview": "node.modify is a unified 'update any node' action. Pass nodeId + any combination of properties. Supports: x, y, width, height, color/fillColor, opacity, cornerRadius, visible, name, rotation, characters/text, fontSize, fontFamily, fontStyle, textAlignHorizontal, layoutSizingHorizontal, layoutSizingVertical, isMask, blendMode.",
				"example":   `{"command":"modify","params":{"nodeId":"$card","fillColor":"#FF0000","cornerRadius":16,"opacity":0.9}}`,
				"tip":       "fillColor is an alias for color — both work. Use modify instead of separate move + resize + paint calls when changing multiple properties at once.",
			},
			"findNodes": map[string]interface{}{
				"_overview": "document.find_nodes is a unified search. Filter by name (query), type (nodeType), and text content (textContent). Returns up to 100 matches.",
				"example":   `{"command":"find","params":{"query":"button","type":"FRAME"}}`,
			},
			"masking": map[string]interface{}{
				"_overview": "node.set_mask creates a mask group. Pass nodeId (the mask shape) and targetIds (nodes to mask). The mask shape becomes isMask=true as the first child of a new group.",
				"usage": []string{
					"Create the mask shape (rect/ellipse/vector) first",
					"Create or reference the target nodes (images, frames)",
					`Call: {"command":"mask","params":{"nodeId":"$ellipse","targetIds":["$image"],"name":"Avatar Mask"}}`,
				},
				"useCases": "Circular avatar crops, shaped image reveals, gradient fade-outs, rounded image cards.",
			},
			"glass": map[string]interface{}{
				"_overview": "effect.apply_glass applies a complete glass morphism recipe in one call: semi-transparent fill + background blur + subtle stroke.",
				"intensities": map[string]interface{}{
					"light":  "8% fill opacity, 20px blur, 10% stroke",
					"medium": "10% fill opacity, 30px blur, 12% stroke",
					"heavy":  "15% fill opacity, 40px blur, 15% stroke",
				},
				"example": `{"command":"glass","params":{"nodeId":"$card","intensity":"medium","tint":"#FFFFFF"}}`,
			},
			"noise": map[string]interface{}{
				"_overview": "effect.add_noise adds a noise overlay effect (Figma Beta API). Creates organic texture on any surface.",
				"types":     "monotone (single color), duotone (two colors), multitone (multiple)",
				"example":   `{"command":"noise","params":{"nodeId":"$bg","noiseType":"monotone","color":"#FFFFFF","noiseSize":100,"density":0.3,"blendMode":"SOFT_LIGHT"}}`,
				"tip":       "Use low density (0.1-0.3) and SOFT_LIGHT blend mode for subtle organic feel. Higher density for grungy/textured looks.",
			},
			"shadowRecipes": map[string]interface{}{
				"_overview": "Production-grade shadow presets. Layer 2-3 shadows for realistic depth.",
				"subtle":    map[string]interface{}{"offsetY": 2, "radius": 4, "spread": 0, "color": "#00000010"},
				"card":      map[string]interface{}{"offsetY": 4, "radius": 12, "spread": -2, "color": "#0000001A"},
				"elevated":  map[string]interface{}{"offsetY": 8, "radius": 24, "spread": -4, "color": "#00000026"},
				"floating":  map[string]interface{}{"offsetY": 16, "radius": 48, "spread": -8, "color": "#00000033"},
				"glow":      map[string]interface{}{"offsetY": 0, "radius": 24, "spread": 4, "color": "accent+40"},
				"innerLight": map[string]interface{}{"type": "INNER_SHADOW", "offsetY": -1, "radius": 0, "spread": 0, "color": "#FFFFFF20"},
				"innerDepth": map[string]interface{}{"type": "INNER_SHADOW", "offsetY": 2, "radius": 4, "spread": 0, "color": "#00000020"},
				"layered": "Combine 2-3 shadows: ambient (large radius, low opacity) + direct (medium) + contact (small radius, close offset). Example: shadow Y:1/R:2/#0D + shadow Y:4/R:12/#1A + shadow Y:16/R:48/#12.",
			},
			"svgIcons": map[string]interface{}{
				"_overview": "Use shape.create_from_svg with inline SVG markup for icons. Common icon libraries: Lucide, Heroicons, Phosphor — LLM can generate the SVG inline.",
				"example":   `{"command":"shape.create_from_svg","params":{"svgPath":"<svg>...</svg>","x":40,"y":40,"width":24,"height":24,"parentId":"$card","name":"icon-check"}}`,
				"tip":       "Set width/height to control icon size. Use paint.set_solid to recolor after creation.",
			},
			"gradientOverlays": map[string]interface{}{
				"_overview": "Use gradient overlays on images for text readability. Create a rect over the image with a transparent-to-dark gradient.",
				"example": `[
  {"name":"overlay","command":"rect","params":{"pid":"$hero","w":1080,"h":400,"y":680,"name":"gradient-overlay"}},
  {"command":"gradient","params":{"nodeId":"$overlay","type":"LINEAR","stops":[{"position":0,"color":"#00000000"},{"position":1,"color":"#000000CC"}]}}
]`,
			},
			"blendModes": map[string]interface{}{
				"MULTIPLY":    "Darkening — great for tinting images, dark overlays",
				"SCREEN":      "Lightening — great for light effects, glows",
				"OVERLAY":     "Contrast boost — great for texture overlays on photos",
				"SOFT_LIGHT":  "Subtle contrast — best for noise/texture overlays",
				"COLOR_DODGE": "Bright highlights — dramatic light effects",
			},
		},
		"stepNamingRules": map[string]interface{}{
			"_rule":    "Batch step names MUST be snake_case. Spaces, hyphens, and special characters are auto-sanitized but you should use clean names from the start.",
			"good":     []string{"hero_bg", "card_title", "gradient_overlay", "cta_button"},
			"bad":      []string{"Hero BG", "card-title", "Gradient Overlay!", "CTA Button"},
			"tip":      "Step names are used for interpolation (${{steps.hero_bg.result.id}}). Clean names prevent lookup failures.",
		},
		"semanticNaming": map[string]interface{}{
			"_rule":    "ALWAYS name every layer by its content/role. Never leave Figma defaults like 'Frame 47' or 'Rectangle 12'.",
			"pattern":  "[Role] - [Detail] or just [Role]",
			"examples": []string{"Hero Background", "Card - Feature Name", "CTA Button", "Gradient Overlay", "Badge Label", "Deco Ring"},
		},
		"imageRules": map[string]interface{}{
			"_overview": "Images in Figma are fills on shapes, not standalone nodes. Use shape.create_image for a one-step convenience.",
			"methods": []string{
				"shape.create_image — ONE STEP: creates a rectangle with an image fill from base64 data. Params: imageData (base64), x, y, width, height, parentId, scaleMode, cornerRadius.",
				"paint.set_image — TWO STEP: first create a shape, then apply an image fill. Good for adding images to existing nodes.",
				"paint.set_image_url — URL-based: fetches an image from a URL (must be in manifest allowedDomains).",
			},
			"filePathSupport": map[string]interface{}{
				"_overview": "imageData accepts file paths and URLs — no need to manually base64 encode. The CLI auto-detects, downloads/reads, and encodes.",
				"formats": []string{
					"file:///absolute/path/to/image.png",
					"/absolute/path/to/image.png",
					"~/Pictures/photo.jpg",
					"https://example.com/image.png",
					"http://localhost:8080/photo.jpg",
				},
				"example_cli":   `ai-happy-design command shape.create_image '{"parentId":"0:1","x":0,"y":0,"width":400,"height":300,"imageData":"file:///tmp/hero.png"}'`,
				"example_url":   `ai-happy-design command shape.create_image '{"parentId":"0:1","x":0,"y":0,"width":400,"height":300,"imageData":"https://picsum.photos/800/600"}'`,
				"example_batch": `{"command":"shape.create_image","params":{"parentId":"$root","imageData":"/tmp/card-bg.jpg","x":0,"y":0,"width":1080,"height":600}}`,
				"supported_ext": ".png, .jpg, .jpeg, .webp, .gif, .svg",
				"with_compress": "Add --compress-images to compress after reading/downloading.",
				"url_note":      "HTTP/HTTPS URLs are downloaded by the CLI (up to 50MB, 30s timeout) and converted to base64 data URIs before sending to the plugin. No allowedDomains restriction since the CLI does the fetching.",
			},
			"scaleModes": []string{
				"FILL — fills the entire shape, cropping if needed (most common)",
				"FIT — fits the image inside the shape, may have empty space",
				"CROP — like FILL but uses a specific crop rectangle",
				"TILE — tiles the image across the shape",
			},
			"compression": []string{
				"For large images (>1MB), add --compress-images flag to batch or command.",
				"This uses ImageMagick to compress imageData before sending (opt-in).",
				"Example: ai-happy-design batch ops.json --compress-images",
				"ImageMagick must be installed (brew install imagemagick or apt install imagemagick).",
			},
			"base64Tips": []string{
				"For CLI batch: write operations to a JSON file, not inline. Base64 images can be >1MB which exceeds shell arg limits.",
				"Example: write JSON to /tmp/img-ops.json, then: ai-happy-design batch /tmp/img-ops.json",
			},
		},
		"quickPrompts": []string{
			"For CREATING designs: use auto-layout frames. Create frames with layoutMode, itemSpacing, padding in one call. Batch for multiple elements.",
			"For EDITING existing nodes: use single commands like paint.set_solid or node.resize.",
			"DEFAULT: Use auto-layout (flexbox) for ALL layout. Absolute positioning only for decorative overlays.",
			"CRITICAL: Always pass lineHeightUnit:'PERCENT' when setting line height. Default is PIXELS which causes overflow.",
			"VERIFY: After creating, call layout.check_overlaps to confirm no elements overlap.",
			"PRINT: Pass dpi:300 to design.compute_tokens for print-ready designs.",
			"FONTS: Call text.list_fonts {fontFamily:'Inter'} to discover available fonts and styles.",
			"BATCH EXPORT: export.batch {nodeIds:'id1,id2', format:'PNG', scale:2} to export multiple frames.",
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
		ops := `[{"name":"card","command":"rect","params":{"x":40,"y":40,"w":220,"h":120,"bg":"#2563EB"}},{"name":"title","command":"text","params":{"text":"10x Faster","x":60,"y":70,"pid":"$card","sz":28,"ff":"Inter","lh":115}}]`
		cliExample = fmt.Sprintf("ai-happy-design batch '%s'", ops)
		mcpExample = map[string]interface{}{
			"tool": "bulk",
			"arguments": map[string]interface{}{
				"action":       "execute",
				"operations":   ops,
				"failFast":     false,
				"retries":      1,
				"retryDelayMs": 250,
				"interpolate":  true,
			},
		}
		batchStep = map[string]interface{}{
			"name":    "bulk_execute",
			"command": "bulk.execute",
			"params": map[string]interface{}{
				"operations":   ops,
				"failFast":     false,
				"retries":      1,
				"retryDelayMs": 250,
				"interpolate":  true,
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
	case strings.EqualFold(param, "failFast"):
		return false
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
