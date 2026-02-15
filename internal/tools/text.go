package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterTextTool registers the "text" tool for text node creation and manipulation.
func RegisterTextTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("text",
		mcp.WithDescription("Text operations: create text nodes, change content, font, size, weight, color, alignment, spacing, case, decoration, and more."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create", "set_content", "set_font", "set_size", "set_weight", "set_color", "set_align", "set_spacing", "set_case", "set_decoration", "get_segments", "load_font", "set_style_id", "list_fonts")),
		mcp.WithString("nodeId", mcp.Description("Target text node ID")),
		mcp.WithString("text", mcp.Description("Text content")),
		mcp.WithString("fontFamily", mcp.Description("Font family name (e.g. Inter, Roboto)")),
		mcp.WithString("fontStyle", mcp.Description("Font style (e.g. Regular, Bold, Italic)")),
		mcp.WithNumber("fontSize", mcp.Description("Font size in pixels")),
		mcp.WithString("fontWeight", mcp.Description("Font weight (e.g. 400, 700)")),
		mcp.WithString("color", mcp.Description("Text color as hex string")),
		mcp.WithString("textAlign", mcp.Description("Text alignment: LEFT, CENTER, RIGHT, JUSTIFIED"),
			mcp.Enum("LEFT", "CENTER", "RIGHT", "JUSTIFIED")),
		mcp.WithNumber("letterSpacing", mcp.Description("Letter spacing value")),
		mcp.WithString("letterSpacingUnit", mcp.Description("Letter spacing unit: PIXELS or PERCENT"),
			mcp.Enum("PIXELS", "PERCENT")),
		mcp.WithNumber("lineHeight", mcp.Description("Line height value")),
		mcp.WithNumber("paragraphSpacing", mcp.Description("Paragraph spacing value")),
		mcp.WithString("textCase", mcp.Description("Text case: ORIGINAL, UPPER, LOWER, TITLE"),
			mcp.Enum("ORIGINAL", "UPPER", "LOWER", "TITLE")),
		mcp.WithString("textDecoration", mcp.Description("Text decoration: NONE, UNDERLINE, STRIKETHROUGH"),
			mcp.Enum("NONE", "UNDERLINE", "STRIKETHROUGH")),
		mcp.WithString("styleId", mcp.Description("Text style ID to apply")),
		mcp.WithNumber("x", mcp.Description("X position for new text")),
		mcp.WithNumber("y", mcp.Description("Y position for new text")),
		mcp.WithString("parentId", mcp.Description("Parent node ID")),
		mcp.WithString("name", mcp.Description("Name for the text node")),
		mcp.WithNumber("width", mcp.Description("Text box width for wrapping")),
		mcp.WithNumber("opacity", mcp.Description("Opacity from 0 to 1")),
		mcp.WithString("lineHeightUnit", mcp.Description("Line height unit: PIXELS, PERCENT, or AUTO"),
			mcp.Enum("PIXELS", "PERCENT", "AUTO")),
		mcp.WithString("layoutSizingHorizontal", mcp.Description("Horizontal sizing in auto-layout: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
		mcp.WithString("layoutSizingVertical", mcp.Description("Vertical sizing in auto-layout: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
		mcp.WithNumber("layoutGrow", mcp.Description("Layout grow: 0 or 1")),
		mcp.WithString("layoutAlign", mcp.Description("Layout alignment: STRETCH or INHERIT"),
			mcp.Enum("STRETCH", "INHERIT")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create":
			params := map[string]interface{}{
				"text":       getStringArg(args, "text", "Hello"),
				"x":          getFloat64Arg(args, "x", 0),
				"y":          getFloat64Arg(args, "y", 0),
				"fontFamily": getStringArg(args, "fontFamily", "Inter"),
				"fontStyle":  getStringArg(args, "fontStyle", "Regular"),
				"fontSize":   getFloat64Arg(args, "fontSize", 16),
				"parentId":   getStringArg(args, "parentId", ""),
				"name":       getStringArg(args, "name", ""),
			}
			// Forward optional params only if explicitly provided
			if hasArg(args, "color") {
				params["color"] = getStringArg(args, "color", "")
			}
			if hasArg(args, "width") {
				params["width"] = getFloat64Arg(args, "width", 0)
			}
			if hasArg(args, "textAlign") {
				params["textAlignHorizontal"] = getStringArg(args, "textAlign", "")
			}
			if hasArg(args, "lineHeight") {
				params["lineHeight"] = getFloat64Arg(args, "lineHeight", 0)
			}
			if hasArg(args, "lineHeightUnit") {
				params["lineHeightUnit"] = getStringArg(args, "lineHeightUnit", "")
			}
			if hasArg(args, "letterSpacing") {
				params["letterSpacing"] = getFloat64Arg(args, "letterSpacing", 0)
			}
			if hasArg(args, "letterSpacingUnit") {
				params["letterSpacingUnit"] = getStringArg(args, "letterSpacingUnit", "")
			}
			if hasArg(args, "textCase") {
				params["textCase"] = getStringArg(args, "textCase", "")
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "layoutSizingHorizontal") {
				params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "")
			}
			if hasArg(args, "layoutSizingVertical") {
				params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "")
			}
			if hasArg(args, "layoutGrow") {
				params["layoutGrow"] = getFloat64Arg(args, "layoutGrow", 0)
			}
			if hasArg(args, "layoutAlign") {
				params["layoutAlign"] = getStringArg(args, "layoutAlign", "")
			}
			return sendCommand(commander, "create_text", params)

		case "set_content":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_content", map[string]interface{}{
				"nodeId": nodeId,
				"text":   getStringArg(args, "text", ""),
			})

		case "set_font":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_font_name", map[string]interface{}{
				"nodeId":     nodeId,
				"fontFamily": getStringArg(args, "fontFamily", "Inter"),
				"fontStyle":  getStringArg(args, "fontStyle", "Regular"),
			})

		case "set_size":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_font_size", map[string]interface{}{
				"nodeId":   nodeId,
				"fontSize": getFloat64Arg(args, "fontSize", 16),
			})

		case "set_weight":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_font_weight", map[string]interface{}{
				"nodeId":     nodeId,
				"fontWeight": getStringArg(args, "fontWeight", "400"),
			})

		case "set_color":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_fill_color", map[string]interface{}{
				"nodeId": nodeId,
				"color":  getStringArg(args, "color", "#000000"),
			})

		case "set_align":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_align", map[string]interface{}{
				"nodeId":    nodeId,
				"textAlign": getStringArg(args, "textAlign", "LEFT"),
			})

		case "set_spacing":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_spacing", map[string]interface{}{
				"nodeId":           nodeId,
				"letterSpacing":    getFloat64Arg(args, "letterSpacing", 0),
				"lineHeight":       getFloat64Arg(args, "lineHeight", 0),
				"paragraphSpacing": getFloat64Arg(args, "paragraphSpacing", 0),
			})

		case "set_case":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_case", map[string]interface{}{
				"nodeId":   nodeId,
				"textCase": getStringArg(args, "textCase", "ORIGINAL"),
			})

		case "set_decoration":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_decoration", map[string]interface{}{
				"nodeId":         nodeId,
				"textDecoration": getStringArg(args, "textDecoration", "NONE"),
			})

		case "get_segments":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "get_styled_text_segments", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "load_font":
			return sendCommand(commander, "load_font_async", map[string]interface{}{
				"fontFamily": getStringArg(args, "fontFamily", "Inter"),
				"fontStyle":  getStringArg(args, "fontStyle", "Regular"),
			})

		case "set_style_id":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_text_style_id", map[string]interface{}{
				"nodeId":  nodeId,
				"styleId": getStringArg(args, "styleId", ""),
			})

		case "list_fonts":
			params := map[string]interface{}{}
			if hasArg(args, "fontFamily") {
				params["family"] = getStringArg(args, "fontFamily", "")
			}
			return sendCommand(commander, "list_fonts", params)

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown text action: %s", action)), nil
		}
	})
}
