package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/nerveband/ai-happy-design/internal/figma"
)

// RegisterBooleanTool registers the "boolean" tool for boolean shape operations.
func RegisterBooleanTool(s *server.MCPServer, commander *figma.Commander) {
	tool := mcp.NewTool("boolean",
		mcp.WithDescription("Boolean operations on shapes: union, subtract, intersect, exclude, flatten."),
		mcp.WithString("action", mcp.Required(), mcp.Description("Action to perform"),
			mcp.Enum("union", "subtract", "intersect", "exclude", "flatten")),
		mcp.WithString("nodeIds", mcp.Description("Comma-separated node IDs (for union, subtract, intersect, exclude)")),
		mcp.WithString("nodeId", mcp.Description("Target node ID (for flatten)")),
	)

	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()
		action := getStringArg(args, "action", "")

		switch action {
		case "union":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "boolean_union", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "subtract":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "boolean_subtract", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "intersect":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "boolean_intersect", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "exclude":
			nodeIds, errResult := requireStringArg(args, "nodeIds")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "boolean_exclude", map[string]interface{}{
				"nodeIds": nodeIds,
			})

		case "flatten":
			nodeId, errResult := requireStringArg(args, "nodeId")
			if errResult != nil {
				return errResult, nil
			}
			return sendCommand(commander, "flatten_node", map[string]interface{}{
				"nodeId": nodeId,
			})

		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown boolean action: %s", action)), nil
		}
	})
}
