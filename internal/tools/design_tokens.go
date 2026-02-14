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
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		width := getFloat64Arg(args, "width", 0)
		height := getFloat64Arg(args, "height", 0)

		if width <= 0 || height <= 0 {
			return mcp.NewToolResultError("width and height must be positive numbers"), nil
		}

		tokens := ComputeDesignTokens(width, height)
		out, _ := json.MarshalIndent(tokens, "", "  ")
		return mcp.NewToolResultText(string(out)), nil
	})
}

// round8 rounds a value to the nearest multiple of 8.
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

// ComputeDesignTokens calculates all design sizing values for a given canvas.
func ComputeDesignTokens(w, h float64) map[string]interface{} {
	ratio := h / w

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
	case ratio >= 1.3:
		layoutType = "portrait"
		columns = "1-2 columns"
		topPadPct = 0.12
		layoutCSS = "body { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:12% 6.5%; gap:4%; }"
		layoutTip = "Portrait/story layout. Lots of vertical space — use it. Center content vertically. Big breathing room top and bottom. Cards should be full-width or 2 max."
	case ratio >= 0.9:
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

	// --- Text sizes (based on WIDTH) ---
	hero := max8(w*0.08, 24)
	heading := max8(w*0.04, 16)
	subheading := max8(w*0.03, 14)
	body := max8(w*0.022, 12)
	caption := max8(w*0.016, 10)
	numbers := max8(w*0.055, 20)
	cta := max8(w*0.022, 12)

	// --- Spacing (horizontal = WIDTH, vertical = HEIGHT) ---
	sidePadding := max8(w*0.065, 16)
	contentWidth := int(w) - 2*sidePadding
	cardPadding := max8(w*0.035, 12)
	framePadding := max8(w*0.06, 16)
	itemSpacing := max8(w*0.015, 8)
	cardGap := max8(w*0.022, 8)
	sectionGap := max8(h*0.035, 16)
	topPadding := max8(h*topPadPct, 16)
	cornerRadius := max8(w*0.015, 8)

	// --- Card sizing ---
	// For multi-column: compute card width
	var cardWidths map[string]interface{}
	cw2 := (contentWidth - cardGap) / 2
	cw3 := (contentWidth - 2*cardGap) / 3
	cw4 := (contentWidth - 3*cardGap) / 4
	cardWidths = map[string]interface{}{
		"fullWidth": contentWidth,
		"twoCol":    cw2,
		"threeCol":  cw3,
		"fourCol":   cw4,
	}

	return map[string]interface{}{
		"_summary": map[string]interface{}{
			"canvas":      map[string]interface{}{"width": int(w), "height": int(h)},
			"aspectRatio": math.Round(ratio*100) / 100,
			"layoutType":  layoutType,
			"cssAnalogy":  layoutCSS,
			"columns":     columns,
			"tip":         layoutTip,
		},
		"text": map[string]interface{}{
			"hero":       hero,
			"heading":    heading,
			"subheading": subheading,
			"body":       body,
			"caption":    caption,
			"numbers":    numbers,
			"cta":        cta,
			"_weights": map[string]interface{}{
				"hero": 800, "heading": 700, "subheading": 600,
				"body": 400, "caption": 500, "numbers": 700, "cta": 700,
			},
			"_lineHeights": map[string]interface{}{
				"hero": 110, "heading": 120, "body": 140, "caption": 130,
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
		"textRules": map[string]interface{}{
			"_critical": []string{
				"TEXT WIDTH: Use the 'width' parameter on text.create to set text wrapping width: text.create {text:'...', x:72, y:100, width:" + itoa(contentWidth) + ", fontSize:24}",
				"LINE HEIGHT: ALWAYS pass lineHeightUnit:'PERCENT'. Example: text.set_line_height {nodeId:id, lineHeight:140, lineHeightUnit:'PERCENT'}. Without this, 140 = 140 PIXELS (not 140%), causing massive overflow.",
				"LAYOUT: Default to absolute x/y positioning. Use auto-layout only for badges/buttons that need centered text.",
				"NEWLINES: NEVER use \\n in text. Create separate text nodes for distinct blocks.",
				"FONTS: Google Fonts are available in Figma. Use fontFamily like 'Space Grotesk', 'Poppins', 'DM Sans'.",
			},
		},
	}
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
