package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design-v2/internal/figma"
)

// RegisterLayerTool registers the "layer" tool for layer ordering and grouping.
func RegisterLayerTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("layer",
		mcp.WithDescription("Layer operations: reorder, bring forward/backward, bring to front/back, group, ungroup, insert child."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("set_order", "bring_forward", "send_backward", "bring_to_front", "send_to_back", "group", "ungroup", "insert_child")),
		mcp.WithString("nodeId", mcp.Description("Target node ID")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for group)")),
		mcp.WithString("parentId", mcp.Description("Parent node ID")),
		mcp.WithString("childId", mcp.Description("Child node ID (for insert_child)")),
		mcp.WithNumber("index", mcp.Description("Target index in parent's children")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "set_order":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_layer_order", map[string]interface{}{
				"nodeId": nodeId,
				"index":  getFloat64Arg(args, "index", 0),
			})

		case "bring_forward":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "bring_forward", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "send_backward":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "send_backward", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "bring_to_front":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "bring_to_front", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "send_to_back":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "send_to_back", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "group":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "group_nodes", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "ungroup":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "ungroup_nodes", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "insert_child":
			parentId, errResult := requireStringArg(args, "parentId")
			if errResult != nil {
				return errResult, nil
			}
			childId, errResult := requireStringArg(args, "childId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "insert_child", map[string]interface{}{
				"parentId": parentId,
				"childId":  childId,
				"index":    getFloat64Arg(args, "index", -1),
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown layer action: %s", action)), nil
		}
	})
}
