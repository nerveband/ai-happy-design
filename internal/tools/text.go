package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design-v2/internal/figma"
)

// RegisterTextTool registers the "text" tool for text node creation and manipulation.
func RegisterTextTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("text",
		mcp.WithDescription("Text operations: create text nodes, change content, font, size, weight, color, alignment, spacing, case, decoration, and more."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create", "set_content", "set_font", "set_size", "set_weight", "set_color", "set_align", "set_spacing", "set_case", "set_decoration", "get_segments", "load_font", "set_style_id")),
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
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create":
			return sendCommand(commander, "create_text", map[string]interface{}{
				"text":       getStringArg(args, "text", "Hello"),
				"x":          getFloat64Arg(args, "x", 0),
				"y":          getFloat64Arg(args, "y", 0),
				"fontFamily": getStringArg(args, "fontFamily", "Inter"),
				"fontStyle":  getStringArg(args, "fontStyle", "Regular"),
				"fontSize":   getFloat64Arg(args, "fontSize", 16),
				"parentId":   getStringArg(args, "parentId", ""),
				"name":       getStringArg(args, "name", ""),
			})

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

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown text action: %s", action)), nil
		}
	})
}
