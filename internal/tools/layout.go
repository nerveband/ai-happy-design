package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterLayoutTool registers the "layout" tool for auto-layout and constraint operations.
func RegisterLayoutTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("layout",
		mcp.WithDescription("Layout operations: auto-layout, padding, spacing, alignment, sizing, wrap, and constraints."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("set_auto_layout", "set_padding", "set_spacing", "set_alignment", "set_sizing", "set_wrap", "set_constraints", "check_overlaps")),
		mcp.WithString("nodeId", mcp.Required(), mcp.Description("Target node ID")),
		mcp.WithString("direction", mcp.Description("Layout direction: HORIZONTAL, VERTICAL, NONE"),
			mcp.Enum("HORIZONTAL", "VERTICAL", "NONE")),
		mcp.WithNumber("paddingTop", mcp.Description("Top padding")),
		mcp.WithNumber("paddingRight", mcp.Description("Right padding")),
		mcp.WithNumber("paddingBottom", mcp.Description("Bottom padding")),
		mcp.WithNumber("paddingLeft", mcp.Description("Left padding")),
		mcp.WithNumber("padding", mcp.Description("Uniform padding on all sides")),
		mcp.WithNumber("itemSpacing", mcp.Description("Spacing between items")),
		mcp.WithString("primaryAxisAlign", mcp.Description("Primary axis alignment: MIN, CENTER, MAX, SPACE_BETWEEN"),
			mcp.Enum("MIN", "CENTER", "MAX", "SPACE_BETWEEN")),
		mcp.WithString("counterAxisAlign", mcp.Description("Counter axis alignment: MIN, CENTER, MAX, BASELINE"),
			mcp.Enum("MIN", "CENTER", "MAX", "BASELINE")),
		mcp.WithString("primaryAxisSizing", mcp.Description("Primary axis sizing: FIXED, AUTO"),
			mcp.Enum("FIXED", "AUTO")),
		mcp.WithString("counterAxisSizing", mcp.Description("Counter axis sizing: FIXED, AUTO"),
			mcp.Enum("FIXED", "AUTO")),
		mcp.WithString("layoutWrap", mcp.Description("Layout wrap mode: NO_WRAP, WRAP"),
			mcp.Enum("NO_WRAP", "WRAP")),
		mcp.WithString("constraintHorizontal", mcp.Description("Horizontal constraint: MIN, CENTER, MAX, STRETCH, SCALE"),
			mcp.Enum("MIN", "CENTER", "MAX", "STRETCH", "SCALE")),
		mcp.WithString("constraintVertical", mcp.Description("Vertical constraint: MIN, CENTER, MAX, STRETCH, SCALE"),
			mcp.Enum("MIN", "CENTER", "MAX", "STRETCH", "SCALE")),
		mcp.WithString("layoutSizingHorizontal", mcp.Description("Child horizontal sizing: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
		mcp.WithString("layoutSizingVertical", mcp.Description("Child vertical sizing: FIXED, HUG, or FILL"),
			mcp.Enum("FIXED", "HUG", "FILL")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")
		nodeId := getStringArg(args, "nodeId", "")

		switch action {
		case "set_auto_layout":
			return sendCommand(commander, "set_auto_layout", map[string]interface{}{
				"nodeId":            nodeId,
				"direction":         getStringArg(args, "direction", "HORIZONTAL"),
				"itemSpacing":       getFloat64Arg(args, "itemSpacing", 0),
				"paddingTop":        getFloat64Arg(args, "paddingTop", 0),
				"paddingRight":      getFloat64Arg(args, "paddingRight", 0),
				"paddingBottom":     getFloat64Arg(args, "paddingBottom", 0),
				"paddingLeft":       getFloat64Arg(args, "paddingLeft", 0),
				"primaryAxisAlign":  getStringArg(args, "primaryAxisAlign", "MIN"),
				"counterAxisAlign":  getStringArg(args, "counterAxisAlign", "MIN"),
			})

		case "set_padding":
			params := map[string]interface{}{"nodeId": nodeId}
			if p := getFloat64Arg(args, "padding", -1); p >= 0 {
				params["paddingTop"] = p
				params["paddingRight"] = p
				params["paddingBottom"] = p
				params["paddingLeft"] = p
			} else {
				params["paddingTop"] = getFloat64Arg(args, "paddingTop", 0)
				params["paddingRight"] = getFloat64Arg(args, "paddingRight", 0)
				params["paddingBottom"] = getFloat64Arg(args, "paddingBottom", 0)
				params["paddingLeft"] = getFloat64Arg(args, "paddingLeft", 0)
			}
			return sendCommand(commander, "set_padding", params)

		case "set_spacing":
			return sendCommand(commander, "set_item_spacing", map[string]interface{}{
				"nodeId":      nodeId,
				"itemSpacing": getFloat64Arg(args, "itemSpacing", 0),
			})

		case "set_alignment":
			return sendCommand(commander, "set_layout_align", map[string]interface{}{
				"nodeId":           nodeId,
				"primaryAxisAlign": getStringArg(args, "primaryAxisAlign", "MIN"),
				"counterAxisAlign": getStringArg(args, "counterAxisAlign", "MIN"),
			})

		case "set_sizing":
			params := map[string]interface{}{
				"nodeId":            nodeId,
				"primaryAxisSizing": getStringArg(args, "primaryAxisSizing", "AUTO"),
				"counterAxisSizing": getStringArg(args, "counterAxisSizing", "AUTO"),
			}
			if hasArg(args, "layoutSizingHorizontal") { params["layoutSizingHorizontal"] = getStringArg(args, "layoutSizingHorizontal", "") }
			if hasArg(args, "layoutSizingVertical") { params["layoutSizingVertical"] = getStringArg(args, "layoutSizingVertical", "") }
			return sendCommand(commander, "set_layout_sizing", params)

		case "set_wrap":
			return sendCommand(commander, "set_layout_wrap", map[string]interface{}{
				"nodeId":     nodeId,
				"layoutWrap": getStringArg(args, "layoutWrap", "NO_WRAP"),
			})

		case "set_constraints":
			return sendCommand(commander, "set_constraints", map[string]interface{}{
				"nodeId":     nodeId,
				"horizontal": getStringArg(args, "constraintHorizontal", "MIN"),
				"vertical":   getStringArg(args, "constraintVertical", "MIN"),
			})

		case "check_overlaps":
			return sendCommand(commander, "check_overlaps", map[string]interface{}{
				"nodeId": nodeId,
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown layout action: %s", action)), nil
		}
	})
}
