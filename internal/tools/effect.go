package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterEffectTool registers the "effect" tool for shadows, blurs, and effect styles.
func RegisterEffectTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("effect",
		mcp.WithDescription("Effect operations: set effects, add shadow, add blur, apply effect style, remove effects."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("set_effects", "add_shadow", "add_blur", "apply_style", "remove")),
		mcp.WithString("nodeId", mcp.Required(), mcp.Description("Target node ID")),
		mcp.WithString("effects", mcp.Description("JSON array of effect objects")),
		mcp.WithString("shadowType", mcp.Description("Shadow type: DROP_SHADOW, INNER_SHADOW"),
			mcp.Enum("DROP_SHADOW", "INNER_SHADOW")),
		mcp.WithString("color", mcp.Description("Shadow color as hex string")),
		mcp.WithNumber("offsetX", mcp.Description("Shadow X offset")),
		mcp.WithNumber("offsetY", mcp.Description("Shadow Y offset")),
		mcp.WithNumber("radius", mcp.Description("Blur or shadow radius")),
		mcp.WithNumber("spread", mcp.Description("Shadow spread")),
		mcp.WithString("blurType", mcp.Description("Blur type: LAYER_BLUR, BACKGROUND_BLUR"),
			mcp.Enum("LAYER_BLUR", "BACKGROUND_BLUR")),
		mcp.WithString("styleId", mcp.Description("Effect style ID to apply")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")
		nodeId := getStringArg(args, "nodeId", "")

		switch action {
		case "set_effects":
			return sendCommand(commander, "set_effects", map[string]interface{}{
				"nodeId":  nodeId,
				"effects": getStringArg(args, "effects", "[]"),
			})

		case "add_shadow":
			return sendCommand(commander, "add_shadow", map[string]interface{}{
				"nodeId":     nodeId,
				"shadowType": getStringArg(args, "shadowType", "DROP_SHADOW"),
				"color":      getStringArg(args, "color", "#00000040"),
				"offsetX":    getFloat64Arg(args, "offsetX", 0),
				"offsetY":    getFloat64Arg(args, "offsetY", 4),
				"radius":     getFloat64Arg(args, "radius", 4),
				"spread":     getFloat64Arg(args, "spread", 0),
			})

		case "add_blur":
			return sendCommand(commander, "add_blur", map[string]interface{}{
				"nodeId":   nodeId,
				"blurType": getStringArg(args, "blurType", "LAYER_BLUR"),
				"radius":   getFloat64Arg(args, "radius", 4),
			})

		case "apply_style":
			styleId, errResult := requireStringArg(args, "styleId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_effect_style_id", map[string]interface{}{
				"nodeId":  nodeId,
				"styleId": styleId,
			})

		case "remove":
			return sendCommand(commander, "set_effects", map[string]interface{}{
				"nodeId":  nodeId,
				"effects": "[]",
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown effect action: %s", action)), nil
		}
	})
}
