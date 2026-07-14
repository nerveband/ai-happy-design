package tools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type MCPPrompt struct {
	Name        string
	Description string
	Text        string
}

func MCPPrompts() []MCPPrompt {
	return []MCPPrompt{
		{Name: "design_strategy", Description: "Plan a Figma design using the AHD catalog and schema-first workflow.", Text: "Use ahd_describe and schema resources first. Compute tokens, find free space, create with batch, then verify structure and visuals."},
		{Name: "batch_strategy", Description: "Use batch for high-throughput Figma canvas operations.", Text: "Prefer one batch for multi-node creation or large edits. Use named steps, compact aliases, dry-run validation, and strict quality checks before reporting completion."},
		{Name: "design_system_strategy", Description: "Work inside an existing design system without drifting styles.", Text: "Analyze existing styles, variables, components, and naming before edits. Reuse tokens and components where possible, then verify with focused reads."},
		{Name: "visual_verification_strategy", Description: "Close the loop with screenshots and local artifact inspection.", Text: "After visual writes, run document.screenshot or document.screenshot_selection, inspect the screenshot artifact, apply corrections, and screenshot again until the result is visually sound."},
		{Name: "accessibility_strategy", Description: "Check accessibility-relevant design qualities before finalizing.", Text: "Validate contrast, text size, hierarchy, spacing, and touch target scale. Use lint and visual verification artifacts to catch issues that schema validation cannot see."},
		{Name: "figma_api_guardrails", Description: "Handle current and beta Figma APIs safely.", Text: "Use runtime-guarded commands for Motion, Shaders, Slots, Grid, Noise, and Texture features. If unavailable, report the unsupported-feature error and continue with compatible alternatives."},
	}
}

func GetMCPPrompt(name string) (MCPPrompt, bool) {
	for _, prompt := range MCPPrompts() {
		if prompt.Name == name {
			return prompt, true
		}
	}
	return MCPPrompt{}, false
}

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
		"version": "3.4",
		"discovery": map[string]interface{}{
			"tools":    "ai-happy-design tools --llm --json",
			"actions":  "ai-happy-design actions [domain]",
			"examples": "ai-happy-design examples [category] — categories: slide, carousel, banner, effects, batch, newsletter",
			"batch":    "ai-happy-design batch --help",
		},
		"lintChecks": map[string]interface{}{
			"_overview":                     "Use --lint on batch to auto-validate designs after creation. Also available as: document.lint {nodeId}",
			"overflow":                      "Children extending beyond parent frame bounds",
			"overlap":                       "Sibling elements overlapping each other (in non-auto-layout frames)",
			"absolute_child_non_autolayout": "Child uses layoutPositioning:ABSOLUTE under a non-auto-layout parent. This is likely unintended and should be fixed.",
			"absolute_overflow":             "Absolute-positioned child exceeds its auto-layout parent bounds.",
			"text_too_large":                "Text fontSize exceeds 50% of parent frame height",
			"text_too_small":                "Text fontSize below 12px (unreadable)",
			"default_name":                  "Nodes with Figma default names (Frame 47, Rectangle 12, etc.)",
			"oversized_child":               "Child width/height exceeds parent by more than 10%",
			"qualityGate":                   "Use --strict-quality to fail the run on any lint warning/error and expose summary.qualityGate in output.",
		},
		"workflow": map[string]interface{}{
			"rule":   "CREATING = batch (one JSON array, many steps). EDITING = single command. Always call document.find_free_space before creating root frames.",
			"create": "Build JSON array → 'ai-happy-design batch ops.json --strict-quality'. Use compact refs ($frame, $last) and aliases (frame/rect/text/fill). Batch auto-normalizes LLM output and auto-places frames.",
			"edit":   "Single command: ai-happy-design command paint.set_solid '{\"nodeId\":\"1:2\",\"color\":\"#FF0000\"}'",
			"verify": "layout.audit {nodeId, compact:true} first for measurable overflow/overlap/text-fit evidence; apply one intentional batch, re-audit, then export.image {nodeId, scale:2} once for visual proof. node.get_tree {compact:true} for structural discovery.",
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
				"_hint":            "Think in CSS, translate to Figma commands.",
				"display:flex":     "layoutMode:VERTICAL/HORIZONTAL + itemSpacing + padding + primaryAxisAlignItems + counterAxisAlignItems",
				"background":       "paint.set_solid or paint.set_gradient",
				"border":           "paint.set_stroke {color, width, strokeAlign:INSIDE}",
				"border-radius":    "node.set_corner_radius {radius}",
				"box-shadow":       "effect.add_shadow {shadowType:DROP_SHADOW, offsetX, offsetY, radius, color}",
				"filter:blur":      "effect.add_blur {blurType:LAYER_BLUR, radius}",
				"backdrop-filter":  "effect.add_blur {blurType:BACKGROUND_BLUR, radius}",
				"opacity":          "node.set_opacity {opacity}",
				"font-size/weight": "text.create {fontSize, fontStyle:'Bold'}",
				"z-index":          "layer.bring_to_front / layer.send_to_back",
				"width:100%":       "layoutSizingHorizontal:FILL in auto-layout",
			},
			"visualHierarchy": map[string]interface{}{
				"_rule":     "Control what viewers see first/second/third using size, weight, color, contrast.",
				"primary":   "display/hero tier, weight 800, accent color. The ONE focal point.",
				"secondary": "heading/subheading tier, weight 600-700, white/slightly muted.",
				"tertiary":  "body/caption tier, weight 400-500, muted gray.",
				"ambient":   "Decorative shapes at very low opacity (0.05-0.15). Must NOT compete with content.",
				"contrast":  "Primary = high contrast, secondary = medium, tertiary = low, ambient = near-invisible.",
			},
			"designDecisions": map[string]interface{}{
				"_rule":     "Match visual treatment to element role.",
				"shadows":   "Cards/buttons/modals get shadows. Recipes: subtle(Y:2,R:8,a:0.1), card(Y:4,R:16,a:0.2), elevated(Y:8,R:32,a:0.3), glow(accent,R:24).",
				"gradients": "Hero backgrounds, CTAs, accent overlays. Avoid on body text and every surface. Use LINEAR with 2 stops.",
				"blur":      "Glass: fill(a:0.08) + BACKGROUND_BLUR(R:16) + stroke(a:0.1). LAYER_BLUR for background depth. Never blur readable content.",
				"strokes":   "Card borders(a:0.07), input fields, decorative outlines, dividers(1px). Stroke-only: set stroke + paint.remove_fill.",
			},
			"layerOrganization": map[string]interface{}{
				"naming":     "Name every element by role: 'Hero Section', 'Card - Feature', 'CTA Button'. Never leave 'Frame 47'.",
				"structure":  "Root Frame → Background (BACK) → Content Frame (auto-layout, MIDDLE) → Overlays (FRONT).",
				"layerOrder": "Back-to-front: bg fills → deco shapes → content frame → foreground overlays. Use layer.send_to_back / bring_to_front.",
			},
			"componentRelationships": map[string]interface{}{
				"containment": "Parent frame owns children (use parentId). Related elements share visual properties.",
				"emphasis":    "Most important element is visually distinct — larger size, accent color, heavier weight.",
				"contrast":    "Adjacent elements need enough contrast. Card fill at 4% white on dark bg + 7% white stroke.",
				"proximity":   "8-16px within a group, 24-48px between groups.",
			},
		},
		"playbook": []string{
			"1. DISCOVER: Get catalog. 2. ANALYZE: design_system.analyze for existing styles. 3. TOKENS: compute_tokens {width, height} for all sizing.",
			"4. FIND SPACE: document.find_free_space for x,y. 5. THINK IN CSS: Draft as HTML/CSS, translate to Figma. 6. HIERARCHY: Use modular type scale tiers.",
			"7. CREATE: batch for multi-element, single command for edits. Batch auto-normalizes LLM output.",
			"8. BALANCE: Siblings match in padding, spacing, radius, text sizes. 9. EFFECTS: Shadows on cards, gradients on hero bg, blur for glass.",
			"10. EXPORT: export.image scale=2 and visually verify. Use node.get_tree compact:true for structural discovery.",
		},
		"designPatterns": map[string]interface{}{
			"_overview":     "These patterns produce professional Figma designs. Follow them strictly.",
			"_decisionTree": "DEFAULT: Use auto-layout (flexbox). Think in rows and columns. Figma auto-layout IS CSS flexbox — you already know it. Create frames with layoutMode, itemSpacing, padding, and alignment in a SINGLE create_frame call. For decorative overlays inside auto-layout frames, use layoutPositioning:ABSOLUTE on the child to exempt it from the flow. After creating, run layout.check_overlaps to verify no elements overlap.",
			"coordinateSystem": map[string]interface{}{
				"origin":     "Top-left (0,0). X right, Y down. x/y relative to parent frame.",
				"autoLayout": "In auto-layout, x/y ignored for flow children. Only absolute-positioned children use x/y.",
				"centering":  "x = (parentWidth - elementWidth) / 2, y = (parentHeight - elementHeight) / 2",
				"snap8px":    "Round all values to nearest multiple of 8.",
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
				"rule":       "Use auto-layout for ALL layout containers. Auto-layout IS CSS flexbox — layoutMode VERTICAL/HORIZONTAL, itemSpacing=gap, padding, primaryAxisAlignItems, counterAxisAlignItems.",
				"oneCommand": `node.create_frame {name:"VStack", width:500, layoutMode:"VERTICAL", itemSpacing:24, padding:32, primaryAxisAlign:"CENTER", counterAxisAlign:"CENTER", primaryAxisSizing:"AUTO"}`,
				"sizing":     "HUG=shrink-wrap (buttons/badges), FILL=stretch to parent (content areas), FIXED=exact dimensions (root frames).",
				"text":       "Create text with width param. Plugin auto-sets layoutSizingVertical=HUG. Use layoutSizingHorizontal=FILL for full-width text.",
				"absolute":   "layoutPositioning:ABSOLUTE only for decorative overlays inside auto-layout. Takes child out of flow.",
				"centering":  "primaryAxisAlign=CENTER + counterAxisAlign=CENTER centers children in both axes.",
			},
			"manualPositioning": map[string]interface{}{
				"rule":          "Non-auto-layout frames position children by x/y (relative to parent). Do NOT set layoutPositioning:ABSOLUTE on children of non-auto-layout frames.",
				"audit":         "Before editing manually positioned content, run layout.audit {nodeId, compact:true}. Use its measured deltas and suggested commands instead of guessing or taking repeated screenshots.",
				"colorShortcut": "node.create_frame, shape.create_rectangle, and text.create all accept a 'color' param (hex or {r,g,b}). NOT 'fillColor' or 'backgroundColor'.",
				"hybrid":        "Root frame (manual x/y) → Badge (auto-layout for centering) → Hero text (x/y, width param) → Cards (x/y) → CTA button (auto-layout for centering)",
			},
			"buildCard": map[string]interface{}{
				"description": "A card is a frame with text children. NEVER create a rectangle + separate floating text.",
				"steps":       "1. create_frame → 2. set fill (opacity 0.04) + cornerRadius + stroke (opacity 0.07) → 3. add text children with parentId. Use consistent heights for sibling cards.",
			},
			"fillRules": map[string]interface{}{
				"_critical":  "Visual surfaces (slide bg, cards, buttons) get fills. Structural frames (wrappers, groups) get NO fill.",
				"visual":     "Use 'color' param on create_frame, or paint.set_solid after.",
				"structural": "Use noFill:true on create_frame, or paint.remove_fill {nodeId, index:0} to clear default white.",
			},
			"layoutPatterns": map[string]interface{}{
				"centeredPage":        "Root frame (FIXED, bg fill) → Content wrapper (VERTICAL auto-layout, CENTER/CENTER, no fill) → children stack and center.",
				"twoColumnGrid":       "HORIZONTAL auto-layout parent, child card frames with equal widths.",
				"headerContentFooter": "Root (VERTICAL, FIXED) → Header (FIXED height) → Content (layoutGrow:1) → Footer (FIXED height).",
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
				"_rule":  "ALWAYS call design.compute_tokens {width, height} for exact values. It returns all font sizes, spacing, padding, card widths, and layout type as concrete pixel values.",
				"scale":  "Modular scale: base = W * 0.044, ratio = 1.333. Steps: caption(-1) → body(0) → subheading(+1) → heading(+2) → title(+3) → hero(+4) → display(+5). Text on 4px grid, spacing on 8px grid.",
				"layout": "Square (H/W 0.9-1.1): 2-3 cards/row, vertical stack. Portrait (1.3-2.0): 1-2 cards/row, generous spacing. Landscape (0.4-0.7): 3-4 cards/row, compact vertical.",
				"print":  "Pass dpi:300 to compute_tokens for print. Letter=2550x3300, A4=2480x3508, businessCard=1050x600 (all at 300dpi).",
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
			"balance": map[string]interface{}{
				"_rule":    "Sibling elements MUST match in size, padding, spacing, cornerRadius, and text sizing. Unbalanced siblings look broken.",
				"cards":    "All sibling cards: same width, height (FIXED sizing), padding, itemSpacing, cornerRadius. Formula: cardWidth = (rowWidth - (N-1) * gap) / N.",
				"text":     "Text at the same hierarchy level across siblings: same fontSize, fontWeight, color.",
				"minSizes": "Card height >= 15-20% of canvas height. Content area >= 60% of canvas height.",
			},
			"framePositioning": map[string]interface{}{
				"_rule": "NEVER guess positions. Call document.find_free_space {width, height, gap} before creating root frames. Returns 8px-grid-snapped x,y.",
			},
			"textRules": map[string]interface{}{
				"width":      "Set width on text.create for wrapping. In auto-layout, use layoutSizingHorizontal:FILL.",
				"lineHeight": "CRITICAL: Always pass lineHeightUnit:'PERCENT'. Without it, value is PIXELS. Good values: hero=110%, heading=120%, body=140%.",
				"newlines":   "Never use \\n. Create separate text nodes for each block. Single text node auto-wraps within width.",
			},
			"cardRules": map[string]interface{}{
				"_rule": "Cards are frames (not rectangles) with text children. Use compute_tokens for sizing. All sibling cards: same height, padding, radius.",
			},
			"compositionTips": []string{
				"Call design.compute_tokens first. Use auto-layout by default. Name all elements descriptively.",
				"Before layout edits, run layout.audit. Apply one bounded batch, re-audit, and screenshot only after the audit is clean.",
				"lineHeightUnit:'PERCENT' is required. layoutPositioning:ABSOLUTE only for decorative overlays in auto-layout.",
				"After creating, call layout.check_overlaps. Export with export.batch at scale:2.",
			},
		},
		"compositeCommands": map[string]interface{}{
			"slide":  "Full social slide. Params: canvas (WxH), bg, gradient, elements[] (eyebrow/headline/body/counter/cta/stats).",
			"banner": "Email banner. Params: canvas (WxH), bg, gradient, dividerX, elements[] (headline/subtitle).",
		},
		"advancedFeatures": map[string]interface{}{
			"modify":    "node.modify — unified update: nodeId + any props (x, y, width, height, color, opacity, cornerRadius, rotation, text, fontSize, etc.).",
			"findNodes": "document.find_nodes — search by name (query), type (nodeType), text content (textContent).",
			"masking":   "node.set_mask {nodeId:maskShape, targetIds:[...]} — circular crops, shaped reveals.",
			"glass":     "effect.apply_glass {nodeId, intensity:light/medium/heavy} — one-call glass morphism.",
			"noise":     "effect.add_noise {nodeId, noiseType:monotone, density:0.1-0.3, blendMode:SOFT_LIGHT} — organic texture (Beta API).",
			"svgIcons":  "shape.create_from_svg {svgPath:'<svg>...</svg>', x, y, width, height, parentId} — inline SVG icons.",
		},
		"naming": map[string]interface{}{
			"steps":  "snake_case for batch steps (hero_bg, card_title). Used in ${{steps.NAME.result.id}} interpolation.",
			"layers": "Name every layer by role: 'Hero Background', 'Card - Feature', 'CTA Button'. Never leave 'Frame 47'.",
		},
		"imageRules": map[string]interface{}{
			"_overview":  "Images in Figma are fills on shapes. shape.create_image is a one-step convenience (creates rect + applies image fill).",
			"methods":    "shape.create_image (one-step), paint.set_image (two-step, existing node), paint.set_image_url (URL-based).",
			"imageData":  "Accepts: base64, data URL, file path (/tmp/img.png, ~/img.jpg, file:///path), or HTTP URL. CLI auto-resolves and encodes.",
			"scaleModes": "FILL (crop to fit), FIT (preserve aspect), CROP, TILE.",
			"compress":   "Add --compress-images flag for large images. Write batch JSON to file (not inline) for large base64 payloads.",
		},
		"quickPrompts": []string{
			"Run 'ai-happy-design tools --llm --json' then 'ai-happy-design actions [domain]' to discover commands.",
			"CREATING: batch for multi-element. EDITING: single command. Use auto-layout (flexbox) by default.",
			"Always pass lineHeightUnit:'PERCENT'. Use --strict-quality on first batch run. Export at scale=2.",
		},
		"cliOnlyCommands": map[string]interface{}{
			"extract":   "ai-happy-design extract file.html --width 1080 --height 1350 — parses HTML/CSS into batch JSON.",
			"benchmark": "ai-happy-design benchmark exec ops.json --runs 3 — times batch execution. benchmark pipe --phase-a-ms N for external LLM timing.",
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
Call describe(action="catalog") to discover all available tools and start designing!

CLI discovery sequence:
  ai-happy-design tools --llm --json
  ai-happy-design actions paint
  ai-happy-design batch --help`
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
		"lintChecks":     catalog["lintChecks"],
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
		return "<base64-or-data-url-or-file-path-or-http-url>"
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
