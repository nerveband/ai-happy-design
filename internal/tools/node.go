package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design-v2/internal/figma"
)

// RegisterNodeTool registers the "node" tool for node inspection, creation, and manipulation.
func RegisterNodeTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("node",
		mcp.WithDescription("Node operations: get info, create frames, move, resize, rotate, opacity, blend mode, visibility, lock, rename, delete, clone."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("get_info", "get_tree", "create_frame", "move", "resize", "rotate", "set_opacity", "set_blend_mode", "set_visibility", "set_locked", "rename", "delete", "clone")),
		mcp.WithString("nodeId", mcp.Description("Target node ID")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for batch operations)")),
		mcp.WithString("name", mcp.Description("Node name")),
		mcp.WithNumber("x", mcp.Description("X position")),
		mcp.WithNumber("y", mcp.Description("Y position")),
		mcp.WithNumber("width", mcp.Description("Width")),
		mcp.WithNumber("height", mcp.Description("Height")),
		mcp.WithNumber("rotation", mcp.Description("Rotation in degrees")),
		mcp.WithNumber("opacity", mcp.Description("Opacity from 0 to 1")),
		mcp.WithString("blendMode", mcp.Description("Blend mode (e.g. NORMAL, MULTIPLY, SCREEN)")),
		mcp.WithBoolean("visible", mcp.Description("Visibility")),
		mcp.WithBoolean("locked", mcp.Description("Lock state")),
		mcp.WithString("parentId", mcp.Description("Parent node ID")),
		mcp.WithString("color", mcp.Description("Fill color as hex")),
		mcp.WithNumber("depth", mcp.Description("Tree traversal depth (for get_tree)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "get_info":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "get_node_info", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "get_tree":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "get_node_tree", map[string]interface{}{
				"nodeId": nodeId,
				"depth":  getFloat64Arg(args, "depth", 3),
			})

		case "create_frame":
			return sendCommand(commander, "create_frame", map[string]interface{}{
				"name":     getStringArg(args, "name", "Frame"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"height":   getFloat64Arg(args, "height", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			})

		case "move":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "move_node", map[string]interface{}{
				"nodeId": nodeId,
				"x":      getFloat64Arg(args, "x", 0),
				"y":      getFloat64Arg(args, "y", 0),
			})

		case "resize":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "resize_node", map[string]interface{}{
				"nodeId": nodeId,
				"width":  getFloat64Arg(args, "width", 100),
				"height": getFloat64Arg(args, "height", 100),
			})

		case "rotate":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "rotate_node", map[string]interface{}{
				"nodeId":   nodeId,
				"rotation": getFloat64Arg(args, "rotation", 0),
			})

		case "set_opacity":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_opacity", map[string]interface{}{
				"nodeId":  nodeId,
				"opacity": getFloat64Arg(args, "opacity", 1.0),
			})

		case "set_blend_mode":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_blend_mode", map[string]interface{}{
				"nodeId":    nodeId,
				"blendMode": getStringArg(args, "blendMode", "NORMAL"),
			})

		case "set_visibility":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_visibility", map[string]interface{}{
				"nodeId":  nodeId,
				"visible": getBoolArg(args, "visible", true),
			})

		case "set_locked":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "set_locked", map[string]interface{}{
				"nodeId": nodeId,
				"locked": getBoolArg(args, "locked", false),
			})

		case "rename":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			name, errResult := requireStringArg(args, "name")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "rename_node", map[string]interface{}{
				"nodeId": nodeId,
				"name":   name,
			})

		case "delete":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "delete_node", map[string]interface{}{
				"nodeId": nodeId,
			})

		case "clone":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "clone_node", map[string]interface{}{
				"nodeId": nodeId,
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown node action: %s", action)), nil
		}
	})
}
