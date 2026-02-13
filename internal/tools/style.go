package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design-v2/internal/figma"
)

// RegisterStyleTool registers the "style" tool for creating and applying styles.
func RegisterStyleTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("style",
		mcp.WithDescription("Style operations: create paint/text/effect styles, apply styles, get all styles, remove styles."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("create_paint", "create_text", "create_effect", "apply", "get_all", "remove")),
		mcp.WithString("nodeId", mcp.Description("Target node ID")),
		mcp.WithString("name", mcp.Description("Style name")),
		mcp.WithString("styleId", mcp.Description("Style ID to apply or remove")),
		mcp.WithString("styleType", mcp.Description("Type of style: FILL, STROKE, TEXT, EFFECT"),
			mcp.Enum("FILL", "STROKE", "TEXT", "EFFECT")),
		mcp.WithString("color", mcp.Description("Color as hex string")),
		mcp.WithString("description", mcp.Description("Style description")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "create_paint":
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_paint_style", map[string]interface{}{
				"name":        name,
				"color":       getStringArg(args, "color", "#000000"),
				"description": getStringArg(args, "description", ""),
			})

		case "create_text":
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_text_style", map[string]interface{}{
				"name":        name,
				"description": getStringArg(args, "description", ""),
			})

		case "create_effect":
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "create_effect_style", map[string]interface{}{
				"name":        name,
				"description": getStringArg(args, "description", ""),
			})

		case "apply":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			styleId, errResult := requireStringArg(args, "styleId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "apply_style", map[string]interface{}{
				"nodeId":    nodeId,
				"styleId":   styleId,
				"styleType": getStringArg(args, "styleType", "FILL"),
			})

		case "get_all":
			return sendCommand(commander, "get_styles", map[string]interface{}{})

		case "remove":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "remove_style", map[string]interface{}{
				"nodeId":    nodeId,
				"styleType": getStringArg(args, "styleType", "FILL"),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown style action: %s", action)), nil
		}
	})
}
