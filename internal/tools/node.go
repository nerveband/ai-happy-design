package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
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
		mcp.WithString("stroke", mcp.Description("Stroke color as hex")),
		mcp.WithNumber("strokeWidth", mcp.Description("Stroke width in px")),
		mcp.WithBoolean("noFill", mcp.Description("Remove frame fill on create")),
		mcp.WithNumber("cornerRadius", mcp.Description("Corner radius for frame")),
		mcp.WithNumber("depth", mcp.Description("Tree traversal depth (for get_tree)")),
		mcp.WithString("layoutMode", mcp.Description("Auto-layout direction: HORIZONTAL or VERTICAL"),
			mcp.Enum("HORIZONTAL", "VERTICAL")),
		mcp.WithNumber("itemSpacing", mcp.Description("Gap between children in auto-layout")),
		mcp.WithNumber("padding", mcp.Description("Uniform padding on all sides")),
		mcp.WithNumber("paddingTop", mcp.Description("Top padding")),
		mcp.WithNumber("paddingRight", mcp.Description("Right padding")),
		mcp.WithNumber("paddingBottom", mcp.Description("Bottom padding")),
		mcp.WithNumber("paddingLeft", mcp.Description("Left padding")),
		mcp.WithString("primaryAxisAlign", mcp.Description("Primary axis alignment: MIN, CENTER, MAX, SPACE_BETWEEN"),
			mcp.Enum("MIN", "CENTER", "MAX", "SPACE_BETWEEN")),
		mcp.WithString("counterAxisAlign", mcp.Description("Counter axis alignment: MIN, CENTER, MAX"),
			mcp.Enum("MIN", "CENTER", "MAX")),
		mcp.WithString("primaryAxisSizing", mcp.Description("Primary axis sizing: FIXED or AUTO"),
			mcp.Enum("FIXED", "AUTO")),
		mcp.WithString("counterAxisSizing", mcp.Description("Counter axis sizing: FIXED or AUTO"),
			mcp.Enum("FIXED", "AUTO")),
		mcp.WithString("layoutWrap", mcp.Description("Layout wrap: WRAP or NO_WRAP"),
			mcp.Enum("WRAP", "NO_WRAP")),
		mcp.WithBoolean("clipsContent", mcp.Description("Whether frame clips content")),
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
			params := map[string]interface{}{
				"name":     getStringArg(args, "name", "Frame"),
				"x":        getFloat64Arg(args, "x", 0),
				"y":        getFloat64Arg(args, "y", 0),
				"width":    getFloat64Arg(args, "width", 100),
				"height":   getFloat64Arg(args, "height", 100),
				"parentId": getStringArg(args, "parentId", ""),
				"color":    getStringArg(args, "color", ""),
			}
			if hasArg(args, "stroke") {
				params["stroke"] = getStringArg(args, "stroke", "")
			}
			if hasArg(args, "strokeWidth") {
				params["strokeWidth"] = getFloat64Arg(args, "strokeWidth", 0)
			}
			if hasArg(args, "noFill") {
				params["noFill"] = getBoolArg(args, "noFill", false)
			}
			if hasArg(args, "opacity") {
				params["opacity"] = getFloat64Arg(args, "opacity", 1)
			}
			if hasArg(args, "cornerRadius") {
				params["cornerRadius"] = getFloat64Arg(args, "cornerRadius", 0)
			}
			// Auto-layout params: only forward if explicitly provided
			if hasArg(args, "layoutMode") {
				params["layoutMode"] = getStringArg(args, "layoutMode", "")
			}
			if hasArg(args, "itemSpacing") {
				params["itemSpacing"] = getFloat64Arg(args, "itemSpacing", 0)
			}
			if hasArg(args, "padding") {
				params["padding"] = getFloat64Arg(args, "padding", 0)
			}
			if hasArg(args, "paddingTop") {
				params["paddingTop"] = getFloat64Arg(args, "paddingTop", 0)
			}
			if hasArg(args, "paddingRight") {
				params["paddingRight"] = getFloat64Arg(args, "paddingRight", 0)
			}
			if hasArg(args, "paddingBottom") {
				params["paddingBottom"] = getFloat64Arg(args, "paddingBottom", 0)
			}
			if hasArg(args, "paddingLeft") {
				params["paddingLeft"] = getFloat64Arg(args, "paddingLeft", 0)
			}
			if hasArg(args, "primaryAxisAlign") {
				params["primaryAxisAlign"] = getStringArg(args, "primaryAxisAlign", "")
			}
			if hasArg(args, "counterAxisAlign") {
				params["counterAxisAlign"] = getStringArg(args, "counterAxisAlign", "")
			}
			if hasArg(args, "primaryAxisSizing") {
				params["primaryAxisSizing"] = getStringArg(args, "primaryAxisSizing", "")
			}
			if hasArg(args, "counterAxisSizing") {
				params["counterAxisSizing"] = getStringArg(args, "counterAxisSizing", "")
			}
			if hasArg(args, "layoutWrap") {
				params["layoutWrap"] = getStringArg(args, "layoutWrap", "")
			}
			if hasArg(args, "clipsContent") {
				params["clipsContent"] = getBoolArg(args, "clipsContent", true)
			}
			return sendCommand(commander, "create_frame", params)

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
