package tools

import (
	"context"
	"encoding/json"
	"math"
	"strconv"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterDesignTool registers the "design" tool for computing design tokens.
// This is a pure computation tool — no Figma connection required.
func RegisterDesignTool(s *server.MCPServer) {
	tool := mcp.NewTool("design",
		mcp.WithDescription("Design helper: compute design tokens (font sizes, spacing, padding, layout) for any canvas dimensions. Call this FIRST before creating any design. No Figma connection required."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("compute_tokens")),
		mcp.WithNumber("width", mcp.Required(), mcp.Description("Canvas width in pixels (e.g., 1080)")),
		mcp.WithNumber("height", mcp.Required(), mcp.Description("Canvas height in pixels (e.g., 1920)")),
		mcp.WithNumber("dpi", mcp.Description("DPI for print designs. Default 72 for screen. Use 300 for print-ready output.")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		width := getFloat64Arg(args, "width", 0)
		height := getFloat64Arg(args, "height", 0)

		if width <= 0 || height <= 0 {
			return mcp.NewToolResultError("width and height must be positive numbers"), nil
		}

		dpi := getFloat64Arg(args, "dpi", 0)
		tokens := ComputeDesignTokens(width, height, dpi)
		out, _ := json.MarshalIndent(tokens, "", "  ")
		return mcp.NewToolResultText(string(out)), nil
	})
}

// round8 rounds a value to the nearest multiple of 8 (used for spacing).
func round8(v float64) int {
	return int(math.Round(v/8)) * 8
}

// max8 returns the larger of round8(v) and floor.
func max8(v float64, floor int) int {
	r := round8(v)
	if r < floor {
		return floor
	}
	return r
}

// round4 rounds a value to the nearest multiple of 4 (used for text sizes —
// the 8px grid is too coarse for a usable type scale).
func round4(v float64) int {
	return int(math.Round(v/4)) * 4
}

// max4 returns the larger of round4(v) and floor.
func max4(v float64, floor int) int {
	r := round4(v)
	if r < floor {
		return floor
	}
	return r
}

// ComputeDesignTokens calculates all design sizing values for a given canvas.
//
// Typography uses a modular scale (perfect fourth = 1.333 ratio) from a base
// proportional to the canvas width. Text sizes use a 4px grid; spacing uses 8px.
//
// Social media context: a 1080px canvas renders at ~375px on a phone screen
// (0.347× scaling). Body text at 36px → ~13px on phone — comfortable reading.
// Anything below 24px in the canvas is hard to read on mobile.
func ComputeDesignTokens(w, h, dpi float64) map[string]interface{} {
	ratio := h / w
	const scaleRatio = 1.333 // perfect fourth — moderate contrast, works for social + web

	// --- Layout classification ---
	var layoutType, layoutCSS, layoutTip string
	var columns string
	var topPadPct float64

	switch {
	case ratio > 2.0:
		layoutType = "ultra-tall"
		columns = "1 column only"
		topPadPct = 0.15
		layoutCSS = "body { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:15% 6.5%; gap:5%; }"
		layoutTip = "Extremely tall canvas. Stack everything vertically. Full-width cards. Very generous spacing between sections."
	case ratio >= 1.15:
		layoutType = "portrait"
		columns = "1-2 columns"
		topPadPct = 0.10
		layoutCSS = "body { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:10% 6.5%; gap:4%; }"
		layoutTip = "Portrait layout (4:5, 9:16, etc.). Generous vertical space — use it. Center content vertically with breathing room. Full-width or 2-column max."
	case ratio >= 0.85:
		layoutType = "square"
		columns = "2-3 columns"
		topPadPct = 0.08
		layoutCSS = "body { display:flex; flex-direction:column; align-items:center; padding:8% 6.5%; gap:3.5%; }"
		layoutTip = "Square layout. Centered content stack. Cards in 2-3 column rows. Balanced vertical spacing."
	default:
		layoutType = "landscape"
		columns = "3-4 columns"
		topPadPct = 0.06
		layoutCSS = "body { display:flex; flex-direction:row; align-items:center; padding:6% 6.5%; gap:2%; }"
		layoutTip = "Landscape/widescreen layout. Height is the constraint. Side-by-side layouts. Use HORIZONTAL auto-layout. Keep content compact vertically."
	}

	// --- Modular type scale (perfect fourth from base) ---
	// Base font = 4.4% of width → 1080px → 48px, 600px → 26px
	// Sized for social media: 1080px canvas renders at ~375px on phones (0.35×).
	// Body at 48px → ~17px on phone = comfortable reading size.
	// Caption at 36px → ~12px on phone = minimum readable.
	// Each step multiplies by 1.333; text rounded to 4px grid.
	// --- Base font size ---
	var rawBase float64
	if dpi > 72 {
		// Print mode: compute from physical dimensions
		physWidthIn := w / dpi
		bodyPt := math.Max(10, math.Min(14, physWidthIn*1.4))
		rawBase = bodyPt * dpi / 72.0
	} else {
		// Screen mode: proportional to canvas width
		rawBase = w * 0.044
	}

	caption := max4(rawBase/scaleRatio, 12)              // scale step -1
	body := max4(rawBase, 14)                            // scale step 0 (base)
	subheading := max4(rawBase*scaleRatio, 16)           // scale step +1
	heading := max4(rawBase*math.Pow(scaleRatio, 2), 20) // scale step +2
	title := max4(rawBase*math.Pow(scaleRatio, 3), 24)   // scale step +3
	hero := max4(rawBase*math.Pow(scaleRatio, 4), 32)    // scale step +4
	display := max4(rawBase*math.Pow(scaleRatio, 5), 40) // scale step +5
	numbers := display                                   // hero stats use display size
	cta := body                                          // CTA button text = body size

	// --- Spacing (8px grid) ---
	sidePadding := max8(w*0.065, 16)
	contentWidth := int(w) - 2*sidePadding
	cardPadding := max8(w*0.035, 12)
	framePadding := max8(w*0.06, 16)
	itemSpacing := max8(w*0.015, 8)
	cardGap := max8(w*0.022, 8)
	sectionGap := max8(h*0.035, 16)
	topPadding := max8(h*topPadPct, 16)
	cornerRadius := max8(w*0.015, 8)

	// --- Button / CTA sizing ---
	// Height ~2.5× font, horizontal padding ~1.25× font, rounded to 8px grid.
	ctaHeight := max8(float64(cta)*2.5, 40)
	ctaPaddingH := max8(float64(cta)*1.25, 16)
	ctaCornerPill := ctaHeight / 2
	ctaCornerRounded := max4(float64(ctaHeight)*0.28, 8)

	// --- Card sizing ---
	cw2 := (contentWidth - cardGap) / 2
	cw3 := (contentWidth - 2*cardGap) / 3
	cw4 := (contentWidth - 3*cardGap) / 4
	cardWidths := map[string]interface{}{
		"fullWidth": contentWidth,
		"twoCol":    cw2,
		"threeCol":  cw3,
		"fourCol":   cw4,
	}

	// Shadow token recommendations per tier
	shadowTokens := map[string]interface{}{
		"caption":    map[string]interface{}{"preset": "subtle", "offsetY": 2, "radius": 4, "color": "#00000010"},
		"body":       map[string]interface{}{"preset": "card", "offsetY": 4, "radius": 12, "color": "#0000001A"},
		"subheading": map[string]interface{}{"preset": "card", "offsetY": 4, "radius": 12, "color": "#0000001A"},
		"heading":    map[string]interface{}{"preset": "elevated", "offsetY": 8, "radius": 24, "color": "#00000026"},
		"title":      map[string]interface{}{"preset": "elevated", "offsetY": 8, "radius": 24, "color": "#00000026"},
		"hero":       map[string]interface{}{"preset": "floating", "offsetY": 16, "radius": 48, "color": "#00000033"},
		"display":    map[string]interface{}{"preset": "floating", "offsetY": 16, "radius": 48, "color": "#00000033"},
	}

	// Tier selection guide
	tierGuide := map[string]interface{}{
		"_rule":       "Select text tier based on content length and importance.",
		"shortText":   "1-3 words (e.g. '250+', 'Go Pro') → title or hero tier",
		"mediumText":  "1-2 sentences (e.g. subtitle, tagline) → heading or subheading tier",
		"longText":    "Paragraphs, descriptions → body tier",
		"metadata":    "Timestamps, footnotes, attribution → caption tier",
		"cta":         "Button text → body tier (cta alias)",
		"heroNumbers": "Big statistics, amounts → display tier",
	}

	result := map[string]interface{}{
		"_summary": map[string]interface{}{
			"canvas":      map[string]interface{}{"width": int(w), "height": int(h)},
			"aspectRatio": math.Round(ratio*100) / 100,
			"layoutType":  layoutType,
			"cssAnalogy":  layoutCSS,
			"columns":     columns,
			"tip":         layoutTip,
		},
		"text": map[string]interface{}{
			"display":    display,
			"hero":       hero,
			"title":      title,
			"heading":    heading,
			"subheading": subheading,
			"body":       body,
			"caption":    caption,
			"numbers":    numbers,
			"cta":        cta,
			"_scale": map[string]interface{}{
				"ratio":     scaleRatio,
				"ratioName": "perfect fourth",
				"base":      body,
				"grid":      "4px (text) / 8px (spacing)",
				"steps":     "caption(-1) → body(0) → subheading(+1) → heading(+2) → title(+3) → hero(+4) → display(+5)",
			},
			"_roles": map[string]interface{}{
				"display":    "Hero statistics, large monetary amounts (e.g. '250+', '$25,000')",
				"hero":       "Campaign headlines, page titles, cover slide text",
				"title":      "Slide headlines, single-message slides, section openers",
				"heading":    "Card titles, secondary headings within sections",
				"subheading": "Tertiary headings, labels, emphasized body text",
				"body":       "Body copy, descriptions, supporting text",
				"caption":    "Footnotes, metadata, timestamps, attribution",
				"numbers":    "Alias for display — use for prominent statistics",
				"cta":        "Button text — same as body size",
			},
			"_weights": map[string]interface{}{
				"display": 800, "hero": 800, "title": 700, "heading": 700,
				"subheading": 600, "body": 400, "caption": 500, "numbers": 700, "cta": 700,
			},
			"_lineHeights": map[string]interface{}{
				"display": 100, "hero": 110, "title": 110, "heading": 120,
				"body": 150, "caption": 130,
				"_unit": "PERCENT — you MUST pass lineHeightUnit:'PERCENT' when calling text.set_line_height. Default unit is PIXELS which causes massive text overflow.",
			},
		},
		"spacing": map[string]interface{}{
			"sidePadding":  sidePadding,
			"contentWidth": contentWidth,
			"framePadding": framePadding,
			"cardPadding":  cardPadding,
			"itemSpacing":  itemSpacing,
			"cardGap":      cardGap,
			"sectionGap":   sectionGap,
			"topPadding":   topPadding,
		},
		"button": map[string]interface{}{
			"fontSize":      cta,
			"height":        ctaHeight,
			"paddingH":      ctaPaddingH,
			"cornerPill":    ctaCornerPill,
			"cornerRounded": ctaCornerRounded,
			"_construction": []string{
				"Create rectangle: shape.create_rectangle {parentId, name:'CTA Button BG', x, y, width:" + itoa(ctaPaddingH*2+cta*6) + ", height:" + itoa(ctaHeight) + "}",
				"Set fill color and corner radius (pill=" + itoa(ctaCornerPill) + " or rounded=" + itoa(ctaCornerRounded) + ")",
				"Add text centered inside: text.create {parentId, x:x+" + itoa(ctaPaddingH) + ", y:y+" + itoa((ctaHeight-cta)/2) + ", fontSize:" + itoa(cta) + ", textAlignHorizontal:'CENTER'}",
				"Or use auto-layout frame for truly centered text: create frame, enable auto-layout, set padding and alignment.",
			},
		},
		"cards": map[string]interface{}{
			"widths":       cardWidths,
			"cornerRadius": cornerRadius,
			"_construction": []string{
				"Create frame: node.create_frame {name, parentId, x, y, width, height} using width from 'widths' above",
				"Set cornerRadius=" + itoa(cornerRadius),
				"Add surface fill: {r:1,g:1,b:1,a:0.04} and stroke: {r:1,g:1,b:1,a:0.07}",
				"Place text children inside with explicit x/y relative to card: text.create {parentId:cardId, x:" + itoa(cardPadding) + ", y:" + itoa(cardPadding) + "}",
				"Card height: set a fixed height that fits your content (e.g., 176-200px for title + description).",
			},
		},
		"tierGuide":    tierGuide,
		"shadowTokens": shadowTokens,
		"textRules": map[string]interface{}{
			"_critical": []string{
				"TEXT WIDTH: Use the 'width' parameter on text.create to set text wrapping width: text.create {text:'...', x:" + itoa(sidePadding) + ", y:100, width:" + itoa(contentWidth) + ", fontSize:" + itoa(body) + "}",
				"LINE HEIGHT: ALWAYS pass lineHeightUnit:'PERCENT'. Example: text.set_line_height {nodeId:id, lineHeight:150, lineHeightUnit:'PERCENT'}. Without this, 150 = 150 PIXELS (not 150%), causing massive overflow.",
				"LAYOUT: Use auto-layout frames for structured sections (header, content, footer groups). Set direction:VERTICAL, padding, itemSpacing, and alignment. Children fill or hug as needed. Use absolute x/y only for decorative overlays.",
				"NEWLINES: NEVER use \\n in text. Create separate text nodes for distinct blocks.",
				"FONTS: Google Fonts are available in Figma. Use fontFamily like 'Space Grotesk', 'Poppins', 'DM Sans'.",
				"PHONE READABILITY: A 1080px canvas renders at ~375px on phones (0.35× scale). Text below " + itoa(caption) + "px will be unreadable on mobile.",
			},
		},
		"template": buildStarterTemplate(layoutType, int(w), int(h)),
		"tips": buildTips(
			layoutType,
			int(w),
			int(h),
			hero,
			body,
			caption,
			sidePadding,
			contentWidth,
		),
		"aliases": map[string]interface{}{
			"fontSize": "display|hero|title|heading|subheading|body|caption|numbers|cta",
			"padding":  "side|frame|card",
			"spacing":  "section|card|item|tight",
			"radius":   "card|button|pill",
			"width":    "content",
			"examples": []string{
				`{"sz":"hero"}`,
				`{"padding":"side","itemSpacing":"section"}`,
				`{"r":"card"}`,
				`{"w":"content"}`,
			},
		},
	}

	// Add print metadata if DPI > 72
	if dpi > 72 {
		physWidthIn := w / dpi
		physHeightIn := h / dpi
		bodyPt := math.Max(10, math.Min(14, physWidthIn*1.4))
		result["_print"] = map[string]interface{}{
			"dpi":            dpi,
			"physicalWidth":  math.Round(physWidthIn*100) / 100,
			"physicalHeight": math.Round(physHeightIn*100) / 100,
			"bodyPt":         math.Round(bodyPt*10) / 10,
			"unit":           "inches",
		}
	}

	return result
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

func buildStarterTemplate(layoutType string, canvasW, canvasH int) []interface{} {
	switch layoutType {
	case "landscape":
		return []interface{}{
			map[string]interface{}{
				"name":    "root",
				"command": "frame",
				"params": map[string]interface{}{
					"name": "Design",
					"w":    canvasW,
					"h":    canvasH,
					"bg":   "#0F172A",
				},
			},
			map[string]interface{}{
				"name":    "content",
				"command": "frame",
				"params": map[string]interface{}{
					"name":             "Content",
					"pid":              "$root",
					"w":                canvasW,
					"h":                canvasH,
					"noFill":           true,
					"layoutMode":       "HORIZONTAL",
					"padding":          "side",
					"itemSpacing":      "card",
					"counterAxisAlign": "CENTER",
				},
			},
			map[string]interface{}{
				"name":    "left_column",
				"command": "frame",
				"params": map[string]interface{}{
					"name":             "Left Column",
					"pid":              "$content",
					"w":                int(float64(canvasW) * 0.46),
					"h":                int(float64(canvasH) * 0.72),
					"noFill":           true,
					"layoutMode":       "VERTICAL",
					"itemSpacing":      "card",
					"primaryAxisAlign": "CENTER",
				},
			},
			map[string]interface{}{
				"name":    "hero_title",
				"command": "text",
				"params": map[string]interface{}{
					"name":      "Hero Title",
					"pid":       "$left_column",
					"text":      "YOUR HEADLINE",
					"sz":        "title",
					"fontStyle": "Bold",
					"color":     "#FFFFFF",
					"w":         int(float64(canvasW) * 0.42),
				},
			},
			map[string]interface{}{
				"name":    "body_copy",
				"command": "text",
				"params": map[string]interface{}{
					"name":  "Body Copy",
					"pid":   "$left_column",
					"text":  "Add one to two supporting sentences.",
					"sz":    "body",
					"color": "#CBD5E1",
					"w":     int(float64(canvasW) * 0.42),
					"lh":    150,
				},
			},
			map[string]interface{}{
				"name":    "right_panel",
				"command": "frame",
				"params": map[string]interface{}{
					"name": "Right Panel",
					"pid":  "$content",
					"w":    int(float64(canvasW) * 0.38),
					"h":    int(float64(canvasH) * 0.72),
					"bg":   "#1E293B",
					"r":    "card",
				},
			},
		}
	default:
		return []interface{}{
			map[string]interface{}{
				"name":    "root",
				"command": "frame",
				"params": map[string]interface{}{
					"name": "Design",
					"w":    canvasW,
					"h":    canvasH,
					"bg":   "#111827",
				},
			},
			map[string]interface{}{
				"name":    "content",
				"command": "frame",
				"params": map[string]interface{}{
					"name":             "Content",
					"pid":              "$root",
					"w":                canvasW,
					"h":                canvasH,
					"noFill":           true,
					"layoutMode":       "VERTICAL",
					"padding":          "side",
					"itemSpacing":      "section",
					"primaryAxisAlign": "CENTER",
					"counterAxisAlign": "CENTER",
				},
			},
			map[string]interface{}{
				"name":    "eyebrow",
				"command": "text",
				"params": map[string]interface{}{
					"name":      "Eyebrow",
					"pid":       "$content",
					"text":      "CATEGORY",
					"sz":        "caption",
					"fontStyle": "Medium",
					"color":     "#9CA3AF",
					"w":         "content",
					"textAlign": "CENTER",
				},
			},
			map[string]interface{}{
				"name":    "hero_title",
				"command": "text",
				"params": map[string]interface{}{
					"name":      "Hero Title",
					"pid":       "$content",
					"text":      "YOUR HEADLINE",
					"sz":        "hero",
					"fontStyle": "Bold",
					"color":     "#FFFFFF",
					"w":         "content",
					"textAlign": "CENTER",
					"lh":        110,
				},
			},
			map[string]interface{}{
				"name":    "body_copy",
				"command": "text",
				"params": map[string]interface{}{
					"name":      "Body Copy",
					"pid":       "$content",
					"text":      "Add one to two supporting sentences.",
					"sz":        "body",
					"color":     "#D1D5DB",
					"w":         "content",
					"textAlign": "CENTER",
					"lh":        150,
				},
			},
			map[string]interface{}{
				"name":    "cta_button",
				"command": "frame",
				"params": map[string]interface{}{
					"name":              "CTA Button",
					"pid":               "$content",
					"bg":                "#3B82F6",
					"r":                 "button",
					"layoutMode":        "HORIZONTAL",
					"paddingLeft":       "card",
					"paddingRight":      "card",
					"paddingTop":        "item",
					"paddingBottom":     "item",
					"primaryAxisSizing": "AUTO",
					"counterAxisSizing": "AUTO",
					"primaryAxisAlign":  "CENTER",
					"counterAxisAlign":  "CENTER",
				},
			},
			map[string]interface{}{
				"name":    "cta_text",
				"command": "text",
				"params": map[string]interface{}{
					"name":      "CTA Text",
					"pid":       "$cta_button",
					"text":      "Get Started",
					"sz":        "cta",
					"fontStyle": "Bold",
					"color":     "#FFFFFF",
				},
			},
		}
	}
}

func buildTips(layoutType string, width, height, hero, body, caption, sidePadding, contentWidth int) []string {
	tips := []string{
		"Use semantic names for every layer. Avoid defaults like 'Frame 47'.",
		"Run batch once, read lint warnings, then fix and rerun before export.",
	}
	switch layoutType {
	case "landscape":
		tips = append([]string{
			"Landscape: use split layouts and keep vertical stacking shallow.",
		}, tips...)
	case "square":
		tips = append([]string{
			"Square: keep hierarchy centered and balanced with 2-3 content groups.",
		}, tips...)
	default:
		tips = append([]string{
			"Portrait: prioritize vertical rhythm and generous spacing between sections.",
		}, tips...)
	}

	tips = append(tips,
		"Token sizes: hero="+itoa(hero)+", body="+itoa(body)+", caption(min)="+itoa(caption)+".",
		"Spacing aliases resolve to side="+itoa(sidePadding)+" and content width="+itoa(contentWidth)+" on "+itoa(width)+"x"+itoa(height)+".",
	)
	return tips
}
